package mconn

import (
	"fmt"
	"reflect"
	"testing"
)

func Test_Conn(t *testing.T) {
	mongodburi := GetMongoDBUrl("127.0.0.1", "27017", "", "")
	t1 := Conn(mongodburi)
	fmt.Println(reflect.TypeOf(t1))
	//fmt.Println(Getsrctype(t1))

}
