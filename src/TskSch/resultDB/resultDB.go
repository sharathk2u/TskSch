package resultDB

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
	"runtime/debug"
	"TskSch/mailer"
	"encoding/json"
)

type Schedule struct {
	Id           int
	Name         string
	Task         string
	Hour         int 	//=>24 hr format
	Minute       int
	Second       int
	Day          int
	Week         int
	R            int
	Update       int
	LastModified time.Time
	InsertedOn   time.Time
}

var Insertedtime time.Time
var ModifiedOn time.Time
var task_id int = 1

//INITIALIZER FOR GETTING SEESION FOR RESULTDB
func ResultdbInit(host string) *mgo.Session {
	session, err := mgo.Dial("mongodb://" + host)
	if err != nil {
		fmt.Println("CAN'T CONNECT TO resultDB", err)
		mailer.Mail("GOSERVE: Unable to connect to the DB", "Unable to establish connection with the mongo database\n\n"+err.Error()+"\n\nStack Trace: ------------------\n\n\n"+string(debug.Stack()))
		return nil
	}
	return session

}

//UPDATER FOR RESULTDB
func UpdateResult(session *mgo.Session, Task_id string, Task_name string, Executed bool, TOE string, TTE string, Pid int, Exec_Stat bool, output string, err string) {

	session.SetMode(mgo.Monotonic, true)

	Col := session.DB("TskSch").C("Result")

	SelectFrom := bson.M{"task_id": Task_id}

	ChangeTo := bson.M{"$set": bson.M{"name": Task_name, "executed": Executed, "toe": TOE, "tte": TTE, "pid": Pid, "exec_stat": Exec_Stat, "output": output, "err": err}}

	Err := Col.Update(SelectFrom, ChangeTo)
	if Err != nil {
		fmt.Println("NOT ABLE TO UPDATE TO THE MONGODB", Err)
	}
}

//INSERT INTO SCHEDULE DB 
func InsertSchedule(session *mgo.Session, js []byte) int {

	var taskJs map[string]interface{}
	json.Unmarshal(js, &taskJs)
	
	Insertedtime = time.Now()
	ModifiedOn = Insertedtime

	session.SetMode(mgo.Monotonic, true)

	c := session.DB("TskSch").C("Schedule")

	Task_cmd := taskJs["Cmd"].(string)
	Task_name := taskJs["Name"].(string)

	Week := int(taskJs["Week"].(float64))
	Day := int(taskJs["Day"].(float64))
	Second := int(taskJs["Second"].(float64))
	Minute := int(taskJs["Minute"].(float64))
	Hour := int(taskJs["Hour"].(float64))
	R := int(taskJs["R"].(float64))
	err := c.Insert(&Schedule{Id: task_id, Name: Task_name, Task: Task_cmd, Hour: Hour, Minute: Minute, Second: Second, Day: Day, Week: Week, R: R, Update: 1, LastModified: ModifiedOn, InsertedOn: Insertedtime})
	if err != nil {
		fmt.Println("NOT ABLE TO ADD TO THE MONGODB", err)
	}
	task_id = task_id + 1
	return task_id - 1
}

//UPDATING THE update bit
func UpdateSchedule(session *mgo.Session, id int, updated int) {

	session.SetMode(mgo.Monotonic, true)

	Col := session.DB("TskSch").C("Schedule")

	SelectFrom := bson.M{"id": id}

	ChangeTo := bson.M{"$set": bson.M{"update": updated}}

	Err := Col.Update(SelectFrom, ChangeTo)
	if Err != nil {
		fmt.Println("NOT ABLE TO UPDATE TO THE MONGODB", Err)
	}
}

//UPDATER FOR SCHEDULES
func Update(session *mgo.Session, js []byte, lastmodified time.Time) {

	var taskJs map[string]interface{}
	json.Unmarshal(js, &taskJs)
	
	session.SetMode(mgo.Monotonic, true)

	Col := session.DB("TskSch").C("Schedule")

	id := int(taskJs["Id"].(float64))
	name := taskJs["Name"].(string)
	task := taskJs["Cmd"].(string)
	week := int(taskJs["Week"].(float64))
	day := int(taskJs["Day"].(float64))
	second := int(taskJs["Second"].(float64))
	minute := int(taskJs["Minute"].(float64))
	hour := int(taskJs["Hour"].(float64))
	r := int(taskJs["R"].(float64))

	SelectFrom := bson.M{"id": id}

	ChangeTo := bson.M{"$set": bson.M{"name": name, "task": task, "hour": hour, "minute": minute, "second": second, "day": day, "week": week, "r": r, "update": 1, "lastmodified": lastmodified}}

	Err := Col.Update(SelectFrom, ChangeTo)
	if Err != nil {
		fmt.Println("NOT ABLE TO UPDATE TO THE MONGODB", Err)
	}
}

//FINDING THE CMD BASED ON ID GIVEN BY SCHEDULER
func Find(cmd_id int,host string) string {
	type task struct {
		Task string
		Name string
	}
	result := task{}

	session := ResultdbInit(host)
	defer func() {
		session.Close()
	}()
	session.SetMode(mgo.Monotonic, true)
	Col := session.DB("TskSch").C("Schedule")

	Err := Col.Find(bson.M{"id": cmd_id}).Select(bson.M{"task": 1, "name": 1}).One(&result)
	if Err != nil {
		fmt.Println("CMD ID", cmd_id, " IS NOT ASSIGNED", Err)
		return ""
	} else {
		return result.Task + ":" + result.Name
	}
}

//PING
func Ping(session *mgo.Session) error {
	err := session.Ping()
	return err
}

//RESTART
func Restart(session *mgo.Session) {
	session.Refresh()
}
