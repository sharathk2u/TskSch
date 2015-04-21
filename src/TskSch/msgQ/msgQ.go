package msgQ

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
)

//INITIALIZER FOR MSG QUEUE
func RedisInit() redis.Conn {
	network := "tcp"
	host := ""
	port := "6379"
	address := host + ":" + port
	Conn, err := redis.Dial(network, address)
	if err != nil {
		fmt.Println("CAN'T DIAL", err)
	}
	return Conn
}

//LENGTH OF MSGQ
func Size(Conn redis.Conn) int {
	val, err := redis.Int(Conn.Do("LLEN", "task"))
	if err != nil {
		fmt.Println("CAN'T POP")
	}
	return val
}

//POP THE ID
func Pop(Conn redis.Conn) string {
	val, err := redis.String(Conn.Do("RPOP", "task"))
	if err != nil {
		fmt.Println("CAN'T POP")
	}
	return val
}

//PUSH THE ID IFF UNABLE TO GET ACCESS OF THE COMMAND RESIDING FILE
func Push(Conn redis.Conn, cmd_id string) {
	_, err := Conn.Do("RPUSH", "task", cmd_id)
	if err != nil {
		fmt.Println("CAN'T PUSH IT BACK")
	}
}

//PING
func Ping(Conn redis.Conn) error{
	_, err := Conn.Do("PING")
	return err
}

