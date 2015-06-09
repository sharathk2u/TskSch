package main

import (
	"fmt"
	"TskSch/resultDB"
	"TskSch/msgQ"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func main() {

	session := resultDB.ResultdbInit("sol-serv-a-d1-1.cloudapp.net")
	defer func() {
		session.Close()
	}()
	
	session.SetMode(mgo.Monotonic, true)
	
	Col1 := session.DB("TskSch").C("Schedule")
	_,Err := Col1.RemoveAll(bson.M{})
	if Err != nil {
		fmt.Println("REMOVE ALL NOT WORKED",Err)
	}
	
	Col2 := session.DB("TskSch").C("Result")
	_,Err = Col2.RemoveAll(bson.M{})
	if Err != nil {
		fmt.Println("REMOVE ALL NOT WORKED",Err)
	}

	conn := msgQ.RedisInit("sol-serv-a-d1-1.cloudapp.net","6379")
	_, err := conn.Do("DEL", "task")
	if err != nil {
	        fmt.Println("REMOVING MSGQ NOT WORKED")
	}
	fmt.Println("REMOVED")
}