package main

import (
	"TskSch/msgQ"
	"TskSch/resultDB"
	"TskSch/schedule"
	"bufio"
	"fmt"
	"os"
	"sync"
)

var cmd_id int = 1
func main() {

	//INITIALIZING THE MONGODB
	session := resultDB.ResultdbInit()

	//INITIALIZING THE REDIS DB
	Conn := msgQ.RedisInit()

	//CLOSING ALL THE CONNECTION
	defer func() {
		session.Close()
		Conn.Close()
	}()

	var Wg     sync.WaitGroup
	ScheduleFile := "/home/unbxd/unbxd/src/TskSch/schedule.txt"
	fd, Err := os.Open(ScheduleFile)
	if Err != nil {
		fmt.Println("Error opening schedule.txt file", Err)
		os.Exit(1)
	}

	reader := bufio.NewReader(fd)
	line, err := Readln(reader)
	for err == nil {
		Wg.Add(1)
		go schedule.Push(&Wg,line,cmd_id,ScheduleFile,session,Conn)
		line, err = Readln(reader)
		cmd_id = cmd_id + 1
	}
	Wg.Wait()
}

//READ LINE BY LINE FROM schedule.txt
func Readln(r *bufio.Reader) (string, error) {
	var (
		isRead   bool  = true
		Err      error = nil
		line, ln []byte
	)
	for isRead && Err == nil {
		line, isRead, Err = r.ReadLine()
		ln = append(ln, line...)
	}
	return string(ln), Err
}
