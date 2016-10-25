package main

import (
	"emtool/common/mconn"
	//	"emtool/common/mlog"

	"encoding/json"
	"flag"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net"
	"os"
	//"reflect"
	"strings"
	"sync"
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

var mongoinfo *MongoInfo

var mutex sync.Mutex
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
	//deal with c like: map[ts:6344943641508708353 t:3 h:-7397456082642233555 v:2 op:c ns:b.$cmd o:map[drop:a]]
	case "i":
		a1 := time.Now().UnixNano()

		err_i := mongoinfo.destClient.DB(dbcoll[0]).C(dbcoll[1]).Insert(oplog["o"])
		a2 := time.Now().UnixNano()
		logger.Println("111: ", a2-a1)
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
func (mongoinfo *MongoInfo) StartReadOplog() {
	oplogDB := mongoinfo.srcclient.DB("local").C("oplog.rs")
	var lastTs bson.MongoTimestamp
	var result bson.M

	var tmp1 int64
	tmp1 = mongoinfo.startpos << 32
	var mongostartts = bson.MongoTimestamp(tmp1)
	oplogquery := bson.M{"ts": bson.M{"$gte": mongostartts}}
	oplogIter := oplogDB.Find(oplogquery).LogReplay().Sort("$natural").Tail(5 * time.Second)

	for {
		for oplogIter.Next(&result) {
			/*
				conn, err := CreateTcpClient()
				if err != nil {
					logger.Println("Connect TCP Server Failed", err.Error())
				}*/

			a1 := time.Now().UnixNano()

			timestamp := result["ts"].(bson.MongoTimestamp) >> 32
			mongoinfo.ApplyOplog(result, result["ns"].(string))
			a2 := time.Now().UnixNano()
			logger.Println("写db耗时: ", a2-a1)
			logger.Println("from: ", mongoinfo.fromhost, " to: ", mongoinfo.tohost, " MongoTimestamp: ", result["ts"], " Unixtimestamp: ", timestamp)

			//	err = SendOplog(conn, result)
			//	CheckErr(err)

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
func StartApplyOplog(result bson.M) {

	timestamp := result["ts"].(bson.MongoTimestamp) >> 32
	mongoinfo.ApplyOplog(result, result["ns"].(string))

	logger.Println("from: ", mongoinfo.fromhost, " to: ", mongoinfo.tohost, " MongoTimestamp: ", result["ts"], " Unixtimestamp: ", timestamp)
}

//Create Tcp Client Connection
func CreateTcpClient() (net.Conn, error) {
	conn, err := net.Dial("tcp", "127.0.0.1:23333")

	return conn, err
}

//Create TCP Server and ReceiveOplog
func CreateTcpServer() {
	//var oplog bson.M
	logger.Println("Start Create TCP Server")
	ln, err := net.Listen("tcp", "127.0.0.1:23333")
	if err != nil {
		logger.Println(err.Error())
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			logger.Println(err.Error())
			continue
		}
		//此处有疑问，是否是顺序写
		ch := make(chan interface{})

		go ReceiveOplog(ch, conn)
		<-ch

	}
}

func ReceiveOplog(ch chan interface{}, conn net.Conn) {
	mutex.Lock()
	defer mutex.Unlock()
	var dt []byte

	var dt1 bson.M
	dec := json.NewDecoder(conn)
	err := dec.Decode(&dt)
	if err != nil {
		logger.Println(err.Error())
	}

	err = bson.Unmarshal(dt, &dt1)
	if err != nil {
		logger.Println(err.Error())
	}

	ch <- dt1

	logger.Println(dt1)
	StartApplyOplog(dt1)

}

func SendOplog(conn net.Conn, result bson.M) error {
	logger.Println(result)
	//cuz json.Marshal will convert big long number to float64, so use bson.Marshal and bson.Unmarshal .
	a, err := bson.Marshal(result)

	if err != nil {
		return err
	}
	encoder := json.NewEncoder(conn)
	err = encoder.Encode(a)

	if err != nil {
		return err
	}
	conn.Close()
	fmt.Println("done")
	return err
}

func CheckErr(err error) {
	if err != nil {
		logger.Println(err.Error())
	}
}

func main() {
	var fromhost, tohost, fromport, toport string
	var startpos int64

	//	var tcpconn net.Conn
	flag.StringVar(&fromhost, "fromhost", "", "the source host")
	flag.StringVar(&tohost, "tohost", "", "the dest host")
	flag.StringVar(&fromport, "fromport", "27017", "the source port")
	flag.StringVar(&toport, "toport", "27017", "the dest port")
	flag.Int64Var(&startpos, "startpos", 0, "the start timestamp")

	flag.Parse()
	//tcpconn = ConnTcp()

	logger = log.New(os.Stdout, "\r\n", log.Ldate|log.Ltime|log.Lshortfile)
	mongoinfo = &MongoInfo{fromhost, tohost, fromport, toport, "", "", 0, nil, nil}
	if fromhost == "" || tohost == "" {
		logger.Println("Please use -help to check the usage")
	} else {
		mongoinfo.InitMongoinfo()
		//CreateTcpServer()
		mongoinfo.StartReadOplog()
		fmt.Println(mongoinfo.srcurl)
		logger.Println("HI")

	}

}
