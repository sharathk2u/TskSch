package main
import(
	"TskSch/scheduler"
	"TskSch/msgQ"
	"TskSch/resultDB"
)
func main(){

	//INITIALIZING THE MONGODB
	Session := resultDB.ResultdbInit()

	//INITIALIZING THE REDIS DB
	Conn := msgQ.RedisInit()

	//CLOSING ALL THE CONNECTION
	defer func(){
		Session.Close()
		Conn.Close()
	}()
	scheduler.Schedule(Session,Conn)
}
