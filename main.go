package main

import (
	"emtool/common/mconn"
	"emtool/common/mlogger"
	"flag"
	"fmt"
)

type MongoInfo struct {
	fromhost string
	tohost   string
	fromport string
	toport   string
	srcurl   string
}

func (mongoinfo *MongoInfo) GetOplog() {
	mongoinfo.srcurl = mconn.GetMongoDBUrl(mongoinfo.fromhost, mongoinfo.fromport, "", "")
	mlogger.Logger("i", "The sourceUrl is "+mongoinfo.srcurl)
}

func main() {
	var fromhost, tohost, fromport, toport string
	flag.StringVar(&fromhost, "fromhost", "", "the source host")
	flag.StringVar(&tohost, "tohost", "", "the dest host")
	flag.StringVar(&fromport, "fromport", "27017", "the source port")
	flag.StringVar(&toport, "toport", "27017", "the dest port")

	flag.Parse()
	mongoinfo := &MongoInfo{fromhost, tohost, fromport, toport, ""}
	if fromhost == "" || tohost == "" {
		fmt.Println("Please use -help to check the usage")
	} else {
		mongoinfo.GetOplog()
	}

}
