package msgQ

import (
    "fmt"
    "github.com/garyburd/redigo/redis"
	"runtime/debug"
	"TskSch/mailer"
)

//INITIALIZER FOR MSG QUEUE
func RedisInit(host string ,port string) redis.Conn {
        network := "tcp"
        address := host + ":" + port
        Conn, err := redis.Dial(network, address)
        if err != nil {
                fmt.Println("CAN'T CONNECT TO msgQ", err)
                mailer.Mail("GOSERVE: Unable to connect to the DB", "Unable to establish connection with the redis database\n\n"+ err.Error()+"\n\nStack Trace: ----- ----------------\n\n\n"+string(debug.Stack()))
                return nil
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
