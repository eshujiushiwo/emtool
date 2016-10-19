package main

import (
	"emtool/common/mconn"
	//	"emtool/common/mlog"
	"flag"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"os"
	//"reflect"
	"log"
	"strings"
	"time"
)

type MongoInfo struct {
	fromhost   string
	tohost     string
	fromport   string
	toport     string
	srcurl     string
	desturl    string
	startpos   int64
	srcclient  *mgo.Session
	destClient *mgo.Session
}

var logger *log.Logger

//Init the mongodbinfo
func (mongoinfo *MongoInfo) InitMongoinfo() {
	mongoinfo.srcurl = mconn.GetMongoDBUrl(mongoinfo.fromhost, mongoinfo.fromport, "", "")
	mongoinfo.desturl = mconn.GetMongoDBUrl(mongoinfo.tohost, mongoinfo.toport, "", "")
	mongoinfo.srcclient = mconn.Conn(mongoinfo.srcurl)
	mongoinfo.destClient = mconn.Conn(mongoinfo.desturl)
}

// Split and Apply Oplog
func (mongoinfo *MongoInfo) ApplyOplog(oplog bson.M, coll string) {
	op := oplog["op"]
	dbcoll := strings.SplitN(coll, ".", 2)
	switch op {
	case "i":
		err_i := mongoinfo.destClient.DB(dbcoll[0]).C(dbcoll[1]).Insert(oplog["o"])
		if err_i != nil {

			logger.Println(err_i.Error())
			if strings.Contains(err_i.Error(), "E11000") {
				logger.Println(err_i.Error())
			} else {
				os.Exit(1)
			}
		}
	case "u":
		err_u := mongoinfo.destClient.DB(dbcoll[0]).C(dbcoll[1]).Update(oplog["o2"], oplog["o"])
		if err_u != nil {
			logger.Println(err_u.Error())
			os.Exit(1)
		}
	case "d":
		err_d := mongoinfo.destClient.DB(dbcoll[0]).C(dbcoll[1]).Remove(oplog["o"])
		if err_d != nil {
			logger.Println(err_d.Error())
			os.Exit(1)
		}
	}
}

//Start replicatie
func (mongoinfo *MongoInfo) StartRestore() {
	oplogDB := mongoinfo.srcclient.DB("local").C("oplog.rs")
	var result bson.M
	var tmp1 int64
	var lastTs bson.MongoTimestamp
	tmp1 = mongoinfo.startpos << 32
	var mongostartts = bson.MongoTimestamp(tmp1)
	oplogquery := bson.M{"ts": bson.M{"$gte": mongostartts}}
	oplogIter := oplogDB.Find(oplogquery).LogReplay().Sort("$natural").Tail(5 * time.Second)
	fmt.Println(oplogquery)
	for {
		for oplogIter.Next(&result) {
			fmt.Println(result["ts"])
			lastTs = result["ts"].(bson.MongoTimestamp)
			timestamp := result["ts"].(bson.MongoTimestamp) >> 32
			mongoinfo.ApplyOplog(result, result["ns"].(string))
			logger.Println("from: ", mongoinfo.fromhost, " to: ", mongoinfo.tohost, " MongoTimestamp: ", result["ts"], " Unixtimestamp: ", timestamp)
		}
		if oplogIter.Err() != nil {
			oplogIter.Close()
		}
		if oplogIter.Timeout() {
			continue
		}
		oplogIter = oplogDB.Find(bson.M{"ts": bson.M{"$gte": lastTs}}).LogReplay().Sort("$natural").Tail(5 * time.Second)

	}

}

func main() {
	var fromhost, tohost, fromport, toport string
	var startpos int64
	flag.StringVar(&fromhost, "fromhost", "", "the source host")
	flag.StringVar(&tohost, "tohost", "", "the dest host")
	flag.StringVar(&fromport, "fromport", "27017", "the source port")
	flag.StringVar(&toport, "toport", "27017", "the dest port")
	flag.Int64Var(&startpos, "startpos", 0, "the start timestamp")

	flag.Parse()

	logger = log.New(os.Stdout, "\r\n", log.Ldate|log.Ltime|log.Lshortfile)
	mongoinfo := &MongoInfo{fromhost, tohost, fromport, toport, "", "", 0, nil, nil}
	if fromhost == "" || tohost == "" {
		logger.Println("Please use -help to check the usage")
	} else {
		mongoinfo.InitMongoinfo()

		mongoinfo.StartRestore()
		fmt.Println(mongoinfo.srcurl)
		logger.Println("HI")

	}

}
