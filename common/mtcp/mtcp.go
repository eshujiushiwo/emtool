package mtcp

import (
	"net"
)

func CreateTcpClient() (net.Conn, error) {
	conn, err := net.Dial("tcp", "127.0.0.1:23333")

	return conn, err
}
func CreateTcpServer() {
	//var oplog bson.M
	logger.Println("Start Create TCP Server")
	ln, err := net.Listen("tcp", "127.0.0.1:23333")
	if err != nil {
		logger.Println(err.Error())z
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
