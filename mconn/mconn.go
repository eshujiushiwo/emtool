package mconn

import (
	"emtool/mlogger"
	"gopkg.in/mgo.v2"
	//"gopkg.in/mgo.v2/bson"
)

/*
type MongoInfo struct {
	Fromhost     string
	Tohost       string
	Fromport     string
	Toport       string
	SrcClient    *mgo.Session
	DestClient   *mgo.Session
	SrcDBConn    *mgo.Database
	DestDBConn   *mgo.Database
	SrcCollConn  *mgo.Collection
	DestCollConn *mgo.Collection
}

*/
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
