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
	Hour int  //=>24 hr format
	Minute int
	Second int
	Day int
	Week int
	R int
	Update int
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
var task_id int = 1
func InsertSchedule(session *mgo.Session , taskJs interface{}){

	Insertedtime = time.Now()
	ModifiedOn = Insertedtime

	session.SetMode(mgo.Monotonic, true) 

	c := session.DB("TskSch").C("Schedule") 

	Task_cmd := taskJs.(map[string]interface{})["cmd"].(string)

	Week := taskJs.(map[string]interface{})["week"].(int)
	Day := taskJs.(map[string]interface{})["day"].(int)
	Second := taskJs.(map[string]interface{})["second"].(int)
	Minute := taskJs.(map[string]interface{})["minute"].(int)
	Hour := taskJs.(map[string]interface{})["hour"].(int)
	R := taskJs.(map[string]interface{})["r"].(int)

	err := c.Insert(&Schedule{Id: task_id , Task: Task_cmd, Hour : Hour , Minute : Minute ,Second : Second , Day : Day , Week : Week ,R : R, Update: 1 , LastModified : ModifiedOn , InsertedOn : Insertedtime})
	if err != nil {
		fmt.Println("NOT ABLE TO ADD TO THE MONGODB",err)
	}
	task_id = task_id +1
}

//INSERT INTO SCHEDULE DB
func IInsertSchedule(session *mgo.Session){
	Insertedtime = time.Now()
	ModifiedOn = Insertedtime
	session.SetMode(mgo.Monotonic, true) 
	c := session.DB("TskSch").C("Schedule") 
	err := c.Insert(
		&Schedule{Id: 1 , Task:"cat ~/unbxd/src/TskSch/command/command.go |wc -l", Hour : 0 , Minute : 0, Second : 20 , Day : 1 ,Week : -1, R : 1 ,Update: 1, LastModified : ModifiedOn , InsertedOn : Insertedtime },
		&Schedule{Id: 2 , Task:"cat ~/unbxd/src/TskSch/add.go | wc -l", Hour : 0 , Minute : 0, Second : 10 , Day : 1 ,Week : -1, R : 1 ,Update : 1, LastModified : ModifiedOn , InsertedOn : Insertedtime },
		&Schedule{Id: 3 , Task:"cat ~/unbxd/src/TskSch/command/command.go |wc -w", Hour : 0 , Minute : 0, Second : 10 , Day : 1 ,Week : -1, R : 1 ,Update : 1, LastModified : ModifiedOn , InsertedOn : Insertedtime },
		&Schedule{Id: 4 , Task:"cat ~/unbxd/src/TskSch/add.go | wc ", Hour : 0 , Minute : 0, Second : 20 , Day : 1 ,Week : -1, R : 1 ,Update : 1, LastModified : ModifiedOn , InsertedOn : Insertedtime },
		&Schedule{Id: 5 , Task:"cat ~/unbxd/src/TskSch/add.go ", Hour : 0 , Minute : 0, Second : 20 , Day : 1 ,Week : -1, R : 1 ,Update : 1, LastModified : ModifiedOn , InsertedOn : Insertedtime },
		&Schedule{Id: 6 , Task:"ct ~/unbxd/src/TskSch/command/command.go |wc ", Hour : 0 , Minute : 0, Second : 40 , Day : 1 ,Week : -1 , R : 1 ,Update : 1, LastModified : ModifiedOn , InsertedOn : Insertedtime },
		&Schedule{Id: 7 , Task:"cat ~/unbxd/src/TskSch/add.go | grep main", Hour : 0 , Minute : 0, Second : 30 , Day : 1 ,Week : -1 , R : 1 ,Update : 1, LastModified : ModifiedOn , InsertedOn : Insertedtime },
		&Schedule{Id: 8 , Task:"cat ~/unbxd/src/TskSch/add.go | wc -w", Hour : 0 , Minute : 0, Second : 30 , Day : 1 ,Week : -1, R : 1 ,Update : 1, LastModified : ModifiedOn , InsertedOn : Insertedtime },
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

	ChangeTo := bson.M{"$set": bson.M{"update" : updated }}

	Err := Col.Update(SelectFrom, ChangeTo)
	if Err != nil {
		fmt.Println("NOT ABLE TO UPDATE TO THE MONGODB", Err)
	}
}

//UPDATE 
func Update(session *mgo.Session,taskJs interface{},lastmodified time.Time){
	session.SetMode(mgo.Monotonic, true)

	Col := session.DB("TskSch").C("Schedule")

	id := taskJs.(map[string]interface{})["id"].(int)
	task := taskJs.(map[string]interface{})["cmd"].(string)
	week := taskJs.(map[string]interface{})["week"].(int)
	day := taskJs.(map[string]interface{})["day"].(int)
	second := taskJs.(map[string]interface{})["second"].(int)
	minute := taskJs.(map[string]interface{})["minute"].(int)
	hour := taskJs.(map[string]interface{})["hour"].(int)
	r := taskJs.(map[string]interface{})["r"].(int)

	SelectFrom := bson.M{"id": id }

	ChangeTo := bson.M{"$set": bson.M{"task" : task , "hour" : hour , "minute" : minute, "second" : second , "day" : day ,"week" : week, "r" : r ,"update" : 1 , "lastmodified" : lastmodified}}

	Err := Col.Update(SelectFrom, ChangeTo)
	if Err != nil {
		fmt.Println("NOT ABLE TO UPDATE TO THE MONGODB", Err)
	}
}

//UPDATE 
func UUpdate(session *mgo.Session,id int,task string,hour int, minute int,second int,day int ,week int ,r int ,updated int,lastmodified time.Time){
	session.SetMode(mgo.Monotonic, true)

	Col := session.DB("TskSch").C("Schedule")

	SelectFrom := bson.M{"id": id }

	ChangeTo := bson.M{"$set": bson.M{"task" : task , "hour" : hour , "minute" : minute, "second" : second , "day" : day ,"week" : week, "r" : r ,"update" : updated , "lastmodified" : lastmodified}}

	Err := Col.Update(SelectFrom, ChangeTo)
	if Err != nil {
		fmt.Println("NOT ABLE TO UPDATE TO THE MONGODB", Err)
	}
}


//FINDING THE CMD BASED ON ID GIVEN BY SCHEDULER
func Find(cmd_id int) string {
	type task struct{
		Task string
	}
	result := task{}

	session := ResultdbInit()
	defer func(){
		session.Close()
	}()
	session.SetMode(mgo.Monotonic, true)
	Col := session.DB("TskSch").C("Schedule")

	Err := Col.Find(bson.M{"id": cmd_id}).Select(bson.M{"task":1}).One(&result)
	if Err != nil {
		fmt.Println("CMD ID",cmd_id," IS NOT ASSIGNED",Err)
		return ""
	}else{
		return result.Task
	}
}
