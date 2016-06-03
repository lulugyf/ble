package main

import (
	"fmt"
	"gyf/ble/q"
	"gyf/ble/q/myzk"
	"log"
	"net"
)

func main() {
	qmap := q.NewQMap()
	defer qmap.Close()

	bleid := "10000001"
	zkaddr := "127.0.0.1:2181"
	mysqladdr := "root:@tcp(127.0.0.1:3306)/idmm2?charset=utf8mb4,utf8&collation=utf8_general_ci"
	version := "8"
	zk := myzk.Connect(zkaddr)
	if zk == nil {
		log.Fatal("can not connect to zookeeper")
	}
	defer zk.Close()

	version = zk.ReadData("/idmm2/configServer/version")

	cfg := q.NewCfg(mysqladdr, version)
	cfg.InitTopics(bleid, qmap)

	port, err := cfg.GetBLEPort(bleid)
	if err != nil {
		log.Fatal("can not get listen_port")
	}

	//	p := q.NewQueue("TRecOprCnttDest", "Sub119Opr")
	//	qmap.Put(p)

	zk.CreateNodeTmp("/idmm2/ble/id.10000001",
		fmt.Sprintf("192.168.56.1:%d jolokia-port:15678", port))

	fmt.Println("start")
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		// handle error
	}
	for {
		conn, err := ln.Accept() // this blocks until connection or error
		if err != nil {
			// handle error
			continue
		}
		go func(conn net.Conn) {
			q.HandleConnection(conn, qmap)
		}(conn) // a goroutine handles conn so that the loop can accept other connections
	}
}
