package main
import(
	"fmt"
	"TskSch/resultDB"
)
func main(){
	session := resultDB.ResultdbInit("","","")
	resultDB.IInsertSchedule(session)
	fmt.Println("INSERTED")
}
