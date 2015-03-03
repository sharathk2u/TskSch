package command

import (
	"TskSch/msgQ"
	"bufio"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"os"
)

//SEARCHING FOR COMMAND BASED ON THE ID POPED FROM MSG QUEUE
func Search(c redis.Conn, cmd_id string) string {

	fd, Err := os.Open("../command.txt")
	if Err != nil {
		msgQ.Push(c, cmd_id)
		fmt.Println("Error opening command.txt file", Err)
		os.Exit(1)
	}
	reader := bufio.NewReader(fd)
	line, err := Readln(reader)
	for err == nil {
		if string(line[0]) == cmd_id {
			return line
		}
		line, err = Readln(reader)
	}
	return ""
}

//READ LINE BY LINE FROM command.txt
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
