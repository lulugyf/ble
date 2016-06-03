package myzk

/*
一些zookeeper操作的封装， 参考： https://mmcgrana.github.io/2014/05/getting-started-with-zookeeper-and-go.html
*/

import (
	"fmt"
	"log"
	//	"os"
	"strings"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

func must(err error) {
	if err != nil {
		panic(err)
	}
}

type myzk struct {
	*zk.Conn
}

func Connect(zksStr string) *myzk {
	//	zksStr := "127.0.0.1:2181" //os.Getenv("ZOOKEEPER_SERVERS")
	zks := strings.Split(zksStr, ",")
	conn, _, err := zk.Connect(zks, time.Second*2)
	must(err)
	return &myzk{conn}
}

func (z *myzk) ReadData(path string) string {
	data, stat, err := z.Get(path)
	if err != nil {
		log.Printf("can not ReadData: %s, %v", path, err)
		return ""
	}
	log.Printf("get:    %+v %+v\n", string(data), stat)
	return string(data)
}

/*
create a node force, is exists then delete
*/
func (z *myzk) CreateNodeForce(path, data string) bool {
	exists, _, err := z.Exists(path) //state
	if err != nil {
		log.Printf("check path[%s] exists failed: %v\n", path, err)
		return false
	}
	if exists { //存在的话， 先删除
		err = z.Delete(path, -1)
		if err != nil {
			log.Printf("delete path[%s] failed: %v", path, err)
			return false
		}
	}
	flags := int32(0)
	acl := zk.WorldACL(zk.PermAll)

	path1, err := z.Create(path, []byte(data), flags, acl)
	if err != nil {
		log.Printf("create path[%s] failed: %v\n", path, err)
		return false
	} else {
		log.Printf("create path[%s] success\n", path1)
	}

	return true
}

// create a ephemeral node
func (z *myzk) CreateNodeTmp(path, data string) bool {
	flags := int32(zk.FlagEphemeral)
	acl := zk.WorldACL(zk.PermAll)

	path1, err := z.Create(path, []byte(data), flags, acl)
	if err != nil {
		log.Printf("create path[%s] failed: %v\n", path, err)
		return false
	} else {
		log.Printf("create path[%s] success\n", path1)
	}

	return true
}

func main_test() {
	conn := Connect("127.0.0.1:2181")
	defer conn.Close()

	flags := int32(0)
	acl := zk.WorldACL(zk.PermAll)

	path, err := conn.Create("/01", []byte("data"), flags, acl)
	must(err)
	fmt.Printf("create: %+v\n", path)

	data, stat, err := conn.Get("/01")
	must(err)
	fmt.Printf("get:    %+v %+v\n", string(data), stat)

	stat, err = conn.Set("/01", []byte("newdata"), stat.Version)
	must(err)
	fmt.Printf("set:    %+v\n", stat)

	err = conn.Delete("/01", -1)
	must(err)
	fmt.Printf("delete: ok\n")

	exists, stat, err := conn.Exists("/01")
	must(err)
	fmt.Printf("exists: %+v %+v\n", exists, stat)
}
