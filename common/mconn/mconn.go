package mconn

import (
	"emtool/common/mlogger"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

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

func Conn(MongoUri string) *mgo.Session {

	MClient, err := mgo.Dial(MongoUri)
	mlogger.Logger("i", "The source url is "+MongoUri)
	if err != nil {
		mlogger.Logger("e", "Connect to "+MongoUri+" Failed!"+err.Error())
	} else {
		mlogger.Logger("i", "Connect to "+MongoUri+" Successed!")
	}
	return MClient

}

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
