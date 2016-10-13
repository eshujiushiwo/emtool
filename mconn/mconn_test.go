package mconn

import (
	"fmt"
	"testing"
)

func Test_GetMongoDBUrl(t *testing.T) {
	mongodburi := GetMongoDBUrl("127.0.0.2", "28017", "", "")
	fmt.Println(mongodburi)
}

func Test_Conn(t *testing.T) {
	mongodburi := GetMongoDBUrl("127.0.0.1", "27017", "root", "root")
	Conn(mongodburi)
}
