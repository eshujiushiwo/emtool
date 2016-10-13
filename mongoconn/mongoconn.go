package mongoconn

import (
	"emtool/mlogger"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func (mongoi) MongoConn(src, dest string) {
	var err, err1 error
	srcClient, err := mgo.Dial(src)

}

func GetMongoDBUrl(addr, port string) string {
	var mongoDBUrl string
	if port == "no" {
		mongoDBUrl = "mongodb://" + addr
	} else {
		mongoDBUrl = "mongodb://" + addr + ":" + port
	}
	return mongoDBUrl
}
