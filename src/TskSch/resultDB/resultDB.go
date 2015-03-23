package resultDB

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
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

type Schedule struct {
	Id int
	Task string
	Time int
	Day int
	Updated int
	LastModified time.Time
	InsertedOn time.Time
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
func UpdateResult(session *mgo.Session, Task_id string, Executed bool, TOE string, TTE string, Pid int, Exec_Stat bool, output string, err string) {

	session.SetMode(mgo.Monotonic, true)

	Col := session.DB("TskSch").C("Result")

	SelectFrom := bson.M{"task_id": Task_id}

	ChangeTo := bson.M{"$set": bson.M{"executed": Executed, "toe": TOE, "tte": TTE, "pid": Pid, "exec_stat": Exec_Stat, "output": output, "err": err}}

	Err := Col.Update(SelectFrom, ChangeTo)
	if Err != nil {
		fmt.Println("NOT ABLE TO UPDATE TO THE MONGODB", Err)
	}
}

//INSERT INTO RESULTDB
func InsertResult(session *mgo.Session){

	session.SetMode(mgo.Monotonic, true) 
	c := session.DB("TskSch").C("Result") 
	err := c.Insert(
		&Result{Task_id: "1" , Executed:false , TOE:"" , TTE:"" , Pid:0 , Exec_Stat:false},
		&Result{Task_id: "2" , Executed:false , TOE:"" , TTE:"" , Pid:0 , Exec_Stat:false},
		&Result{Task_id: "3" , Executed:false , TOE:"" , TTE:"" , Pid:0 , Exec_Stat:false},
		&Result{Task_id: "4" , Executed:false , TOE:"" , TTE:"" , Pid:0 , Exec_Stat:false},
		&Result{Task_id: "5" , Executed:false , TOE:"" , TTE:"" , Pid:0 , Exec_Stat:false},
		&Result{Task_id: "6" , Executed:false , TOE:"" , TTE:"" , Pid:0 , Exec_Stat:false},
		&Result{Task_id: "7" , Executed:false , TOE:"" , TTE:"" , Pid:0 , Exec_Stat:false},
		&Result{Task_id: "8" , Executed:false , TOE:"" , TTE:"" , Pid:0 , Exec_Stat:false})
	if err != nil {
		fmt.Println("NOT ABLE TO ADD TO THE MONGODB",err)
	}
}

var Insertedtime time.Time
var ModifiedOn time.Time

//INSERT INTO SCHEDULE DB
func InsertSchedule(session *mgo.Session){
	Insertedtime = time.Now()
	ModifiedOn = Insertedtime
	session.SetMode(mgo.Monotonic, true) 
	c := session.DB("TskSch").C("Schedule") 
	err := c.Insert(
		&Schedule{Id: 1 , Task:"cat ~/unbxd/src/TskSch/command/command.go |wc -l", Time : 20 , Day : 0 ,Updated : 1, LastModified : ModifiedOn , InsertedOn : Insertedtime },
		&Schedule{Id: 2 , Task:"cat ~/unbxd/src/TskSch/command.txt | wc -l", Time : 30 , Day : 0 ,Updated : 1, LastModified : ModifiedOn , InsertedOn : Insertedtime },
		&Schedule{Id: 3 , Task:"ct ~/unbxd/src/TskSch/command/command.go |wc -w", Time : 20 , Day : 0 ,Updated : 1, LastModified : ModifiedOn , InsertedOn : Insertedtime },
		&Schedule{Id: 4 , Task:"cat ~/unbxd/src/TskSch/command.txt | wc ", Time : 40 , Day : 0 ,Updated : 1, LastModified : ModifiedOn , InsertedOn : Insertedtime },
		&Schedule{Id: 5 , Task:"cat ~/unbxd/src/TskSch/command.txt ", Time : 10 , Day : 0 ,Updated : 1, LastModified : ModifiedOn , InsertedOn : Insertedtime },
		&Schedule{Id: 6 , Task:"cat ~/unbxd/src/TskSch/command/command.go |wc ", Time : 40 , Day : 0 ,Updated : 1, LastModified : ModifiedOn , InsertedOn : Insertedtime },
		&Schedule{Id: 7 , Task:"cat ~/unbxd/src/TskSch/command.txt |grep cat", Time : 50 , Day : 1 ,Updated : 1, LastModified : ModifiedOn , InsertedOn : Insertedtime },
		&Schedule{Id: 8 , Task:"cat ~/unbxd/src/TskSch/command.txt |wc -w", Time : 10 , Day : 1 ,Updated : 1, LastModified : ModifiedOn , InsertedOn : Insertedtime },
)
	if err != nil {
		fmt.Println("NOT ABLE TO ADD TO THE MONGODB",err)
	}
}

//UPDATE SCHEDULES
func UpdateSchedule(session *mgo.Session, id int,updated int){

	session.SetMode(mgo.Monotonic, true)

	Col := session.DB("TskSch").C("Schedule")

	SelectFrom := bson.M{"id": id }

	ChangeTo := bson.M{"$set": bson.M{"updated" : updated }}

	Err := Col.Update(SelectFrom, ChangeTo)
	if Err != nil {
		fmt.Println("NOT ABLE TO UPDATE TO THE MONGODB", Err)
	}
}

//UPDATE 
func Update(session *mgo.Session,id int,task string,time int, day int,updated int,lastmodified time.Time){
	session.SetMode(mgo.Monotonic, true)

	Col := session.DB("TskSch").C("Schedule")

	SelectFrom := bson.M{"id": id }

	ChangeTo := bson.M{"$set": bson.M{"task" : task , "time" : time , "day" : day,"updated" : updated , "lastmodified" : lastmodified}}

	Err := Col.Update(SelectFrom, ChangeTo)
	if Err != nil {
		fmt.Println("NOT ABLE TO UPDATE TO THE MONGODB", Err)
	}
}
