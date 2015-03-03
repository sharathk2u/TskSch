package resultDB

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Result struct {
	Task_id   string //command ID
	Executed  bool   //Executed staus
	TOE       string //Time Of Execution
	TTE       string //Time Taken to Execute
	Pid       int    // process id of client
	Exec_Stat bool   //Execution Status
	output    string
	err       string
}

//INITIALIZER FOR GETTING SEESION FOR RESULTDB
func ResultdbInit() *mgo.Session {
	host := "localhost"
	session, err := mgo.Dial("mongodb://" + host)
	if err != nil {
		fmt.Println("NOT ABLE TO CONNECT TO MONGODB SERVER", err)
	}
	return session

}

//UPDATER FOR RESULTDB
func Update(session *mgo.Session, Task_id string, Executed bool, TOE string, TTE string, Pid int, Exec_Stat bool, output string, err string) {

	session.SetMode(mgo.Monotonic, true)

	Col := session.DB("TskSch").C("Result")

	SelectFrom := bson.M{"task_id": Task_id}

	ChangeTo := bson.M{"$set": bson.M{"executed": Executed, "toe": TOE, "tte": TTE, "pid": Pid, "exec_stat": Exec_Stat, "output": output, "err": err}}

	Err := Col.Update(SelectFrom, ChangeTo)
	if Err != nil {
		fmt.Println("NOT ABLE TO UPDATE TO THE MONGODB", Err)
	}
}
