package main
import(
	"fmt"
	"TskSch/resultDB"
	"time"
)
func main(){
	session := resultDB.ResultdbInit()
	resultDB.UUpdate(session,1,"ps",0,0,20,1,-1,1,1,time.Now())
	fmt.Println("UPDATED")
}
