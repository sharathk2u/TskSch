package main
import(
        "fmt"
        "TskSch/msgQ"
)
func main(){
        conn := msgQ.RedisInit("sol-serv-a-d1-1.cloudapp.net","6379")
	msgQ.Push(conn , "1")
	fmt.Println("INSERTED")
}
