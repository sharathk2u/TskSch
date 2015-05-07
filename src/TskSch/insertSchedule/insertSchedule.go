package main
import(
        "fmt"
        "TskSch/resultDB"
)
func main(){
        session := resultDB.ResultdbInit("sol-serv-a-d1-1.cloudapp.net")
        resultDB.IInsertSchedule(session)
        fmt.Println("INSERTED")
}
