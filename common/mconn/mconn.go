package mconn

import (
	//"emtool/common/mlog"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"os"
)

//Get Mongodb Connection strings
var logger *log.Logger

func GetMongoDBUrl(addr, port, username, password string) string {
	var mongoDBUrl string
	if port == "no" {

		if username == "" || password == "" {

			mongoDBUrl = "mongodb://" + addr
		} else {
			mongoDBUrl = "mongodb://" + username + ":" + password + "@" + addr
		}

	} else {
		if username == "" || password == "" {

			mongoDBUrl = "mongodb://" + addr + ":" + port
		} else {
			mongoDBUrl = "mongodb://" + username + ":" + password + "@" + addr + ":" + port
		}
	}
	return mongoDBUrl
}

// Build Connections
func Conn(MongoUri string) *mgo.Session {
	MClient, err := mgo.Dial(MongoUri)

	if err != nil {

		logger.Println("Connect to ", MongoUri, " Failed!", err.Error())
		os.Exit(-1)
	} else {
		logger.Println("Connect to ", MongoUri, " Successed!")
	}
	return MClient

}

// Check the node is mongod or not
func Getsrctype(Mclient *mgo.Session) string {
	var srctype string
	command := bson.M{"isMaster": 1}
	result := bson.M{}
	Mclient.Run(command, &result)
	if result["msg"] == "isdbgrid" {
		srctype = "mongos"
	} else {
		srctype = "mongod"
	}
	return srctype

}
func init() {
	logger = log.New(os.Stdout, "\r\n", log.Ldate|log.Ltime|log.Lshortfile)
}
