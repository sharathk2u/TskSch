package main
import(
	"TskSch/scheduler"
	"TskSch/resultDB"
)
func main(){

	//INITIALIZING THE MONGODB
	Session := resultDB.ResultdbInit()

	//CLOSING ALL THE CONNECTION
	defer func(){
		Session.Close()
	}()
	scheduler.Schedule(Session)
}
