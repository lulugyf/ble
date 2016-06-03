package q

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type cfg struct {
	db      *sql.DB
	version string
}

func NewCfg(connstr string, version string) *cfg {
	// "root:@tcp(127.0.0.1:3306)/test?charset=utf8mb4,utf8&collation=utf8_general_ci"
	db, err := sql.Open("mysql", connstr)
	if err != nil {
		log.Printf("connect mysql failed: %v\n", err)
		return nil
	}
	return &cfg{db: db, version: version}
}

func (c *cfg) GetBLEPort(bleid string) (int, error) {
	sql := `select addr_ip, addr_port 
		from ble_base_info_{version}
		where ble_id=? and id_number=?`
	sql = strings.Replace(sql, "{version}", c.version, -1)
	stmt1, err := c.db.Prepare(sql)
	if err != nil {
		log.Printf("prepare sql failed: %v", err)
		return -1, errors.New("prepare statement failed")
	}
	defer stmt1.Close()
	var ip string
	var port int
	err = stmt1.QueryRow(bleid, 1).Scan(&ip, &port)
	if err != nil {
		log.Printf("query port failed: %v", err)
		return -1, errors.New("failed")
	}
	return port, nil
}
func (c *cfg) InitTopics(bleid string, qm *qmap) error {
	sql := `select a.ble_id, b.client_id, a.dest_topic_id, b.max_request, b.min_timeout, b.max_timeout
				from ble_dest_topic_rel_{version} a, topic_subscribe_rel_{version} b 
				where a.use_status='1' and b.use_status='1' 
				and a.dest_topic_id=b.dest_topic_id `
	// and a.ble_id=?
	sql = strings.Replace(sql, "{version}", c.version, -1)
	log.Printf("==sql:%s\n", sql)
	stmt1, err := c.db.Prepare(sql)
	if err != nil {
		log.Printf("prepare sql failed: %v", err)
		return errors.New("prepare statement failed")
	}
	defer stmt1.Close()
	rows, err := stmt1.Query() //bleid
	if err != nil {
		log.Printf("query failed: %v", err)
		return errors.New("query topics failed")
	}
	var ble_id, clientid, topic string
	var max_request, min_timeout, max_timeout int
	for rows.Next() {
		e1 := rows.Scan(&ble_id, &clientid, &topic,
			&max_request, &min_timeout, &max_timeout)
		if e1 != nil {
			break
		}
		log.Printf("%s - %s - %s - %d %d %d", ble_id, clientid, topic,
			max_request, min_timeout, max_timeout)
		if bleid == ble_id && qm != nil {
			p := NewQueue(topic, clientid)
			p.max_request = max_request
			p.min_timeout = min_timeout
			p.max_timeout = max_timeout
			qm.Put(p)
		}

	}
	return nil
}

func Mysql_test() {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/test?charset=utf8mb4,utf8&collation=utf8_general_ci")
	if err != nil {
		log.Printf("connect mysql failed: %v\n", err)
		return
	}
	defer db.Close()

	log.Printf("connect to mysql successfully!\n")

	//create table squarenum(`number` int, `squareNumber` int)
	// Prepare statement for inserting data
	stmtIns, err := db.Prepare("INSERT INTO squareNum VALUES( ?, ? )") // ? = placeholder
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmtIns.Close() // Close the statement when we leave main() / the program terminates

	// Prepare statement for reading data
	stmtOut, err := db.Prepare("SELECT squareNumber FROM squarenum WHERE number = ?")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmtOut.Close()

	// Insert square numbers for 0-24 in the database
	for i := 0; i < 25; i++ {
		_, err = stmtIns.Exec(i, (i * i)) // Insert tuples (i, i^2)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
	}

	var squareNum int // we "scan" the result in here

	// Query the square-number of 13
	err = stmtOut.QueryRow(13).Scan(&squareNum) // WHERE number = 13
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	fmt.Printf("The square number of 13 is: %d", squareNum)

	// Query another number.. 1 maybe?
	err = stmtOut.QueryRow(1).Scan(&squareNum) // WHERE number = 1
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	fmt.Printf("The square number of 1 is: %d", squareNum)
}
