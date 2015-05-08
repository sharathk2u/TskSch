package resultDB

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Result struct {
	Task_id   string //command ID
	Task_name string
	Executed  bool   //Executed staus
	TOE       string //Time Of Execution
	TTE       string //Time Taken to Execute
	Pid       int    // process id of client
	Exec_Stat bool   //Execution Status
	output    string
	err       string
}

type Schedule struct {
	Id           int
	Name         string
	Task         string
	Hour         int //=>24 hr format
	Minute       int
	Second       int
	Day          int
	Week         int
	R            int
	Update       int
	LastModified time.Time
	InsertedOn   time.Time
}

//INITIALIZER FOR GETTING SEESION FOR RESULTDB
func ResultdbInit(host string) *mgo.Session {
	session, err := mgo.Dial("mongodb://" + host)
	if err != nil {
		fmt.Println("NOT ABLE TO CONNECT TO MONGODB SERVER", err)
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

//INSERT INTO RESULTDB
func InsertResult(session *mgo.Session) {

	session.SetMode(mgo.Monotonic, true)
	c := session.DB("TskSch").C("Result")
	err := c.Insert(
		&Result{Task_id: "1", Task_name: "A", Executed: false, TOE: "", TTE: "", Pid: 0, Exec_Stat: false},
		&Result{Task_id: "2", Task_name: "B", Executed: false, TOE: "", TTE: "", Pid: 0, Exec_Stat: false},
		&Result{Task_id: "3", Task_name: "C", Executed: false, TOE: "", TTE: "", Pid: 0, Exec_Stat: false},
		&Result{Task_id: "4", Task_name: "D", Executed: false, TOE: "", TTE: "", Pid: 0, Exec_Stat: false},
		&Result{Task_id: "5", Task_name: "E", Executed: false, TOE: "", TTE: "", Pid: 0, Exec_Stat: false},
		&Result{Task_id: "6", Task_name: "F", Executed: false, TOE: "", TTE: "", Pid: 0, Exec_Stat: false},
		&Result{Task_id: "7", Task_name: "G", Executed: false, TOE: "", TTE: "", Pid: 0, Exec_Stat: false},
		&Result{Task_id: "8", Task_name: "H", Executed: false, TOE: "", TTE: "", Pid: 0, Exec_Stat: false})
	if err != nil {
		fmt.Println("NOT ABLE TO ADD TO THE MONGODB", err)
	}
}

var Insertedtime time.Time
var ModifiedOn time.Time
var task_id int = 1

func InsertSchedule(session *mgo.Session, taskJs map[string]interface{}) int {

	Insertedtime = time.Now()
	ModifiedOn = Insertedtime

	session.SetMode(mgo.Monotonic, true)

	c := session.DB("TskSch").C("Schedule")
	Task_cmd := taskJs["cmd"].(string)
	Task_name := taskJs["name"].(string)

	Week := int(taskJs["week"].(float64))
	Day := int(taskJs["day"].(float64))
	Second := int(taskJs["second"].(float64))
	Minute := int(taskJs["minute"].(float64))
	Hour := int(taskJs["hour"].(float64))
	R := int(taskJs["r"].(float64))

	err := c.Insert(&Schedule{Id: task_id, Name: Task_name, Task: Task_cmd, Hour: Hour, Minute: Minute, Second: Second, Day: Day, Week: Week, R: R, Update: 1, LastModified: ModifiedOn, InsertedOn: Insertedtime})
	if err != nil {
		fmt.Println("NOT ABLE TO ADD TO THE MONGODB", err)
	}
	task_id = task_id + 1
	return task_id - 1
}

//INSERT INTO SCHEDULE DB
func IInsertSchedule(session *mgo.Session) {
	Insertedtime = time.Now()
	ModifiedOn = Insertedtime
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("TskSch").C("Schedule")
	err := c.Insert(
		&Schedule{Id: 1, Name: "A", Task: "ls  | wc -l", Hour: 0, Minute: 0, Second: 20, Day: 1, Week: -1, R: 1, Update: 1, LastModified: ModifiedOn, InsertedOn: Insertedtime},
		&Schedule{Id: 2, Name: "B", Task: "ls -l | wc -l", Hour: 0, Minute: 0, Second: 10, Day: 1, Week: -1, R: 1, Update: 1, LastModified: ModifiedOn, InsertedOn: Insertedtime},
		&Schedule{Id: 3, Name: "C", Task: "ls -la |wc -w", Hour: 0, Minute: 0, Second: 10, Day: 1, Week: -1, R: 1, Update: 1, LastModified: ModifiedOn, InsertedOn: Insertedtime},
		&Schedule{Id: 4, Name: "D", Task: "l -l | wc ", Hour: 0, Minute: 0, Second: 20, Day: 1, Week: -1, R: 1, Update: 1, LastModified: ModifiedOn, InsertedOn: Insertedtime},
		&Schedule{Id: 5, Name: "E", Task: "ls -la", Hour: 0, Minute: 0, Second: 20, Day: 1, Week: -1, R: 1, Update: 1, LastModified: ModifiedOn, InsertedOn: Insertedtime},
		&Schedule{Id: 6, Name: "F", Task: "ls | grep .go | wc -w ", Hour: 0, Minute: 0, Second: 40, Day: 1, Week: -1, R: 1, Update: 1, LastModified: ModifiedOn, InsertedOn: Insertedtime},
		&Schedule{Id: 7, Name: "G", Task: "ls -la | grep .go", Hour: 0, Minute: 0, Second: 30, Day: 1, Week: -1, R: 1, Update: 1, LastModified: ModifiedOn, InsertedOn: Insertedtime},
		&Schedule{Id: 8, Name: "H", Task: "ls -a | wc -l", Hour: 0, Minute: 0, Second: 30, Day: 1, Week: -1, R: 1, Update: 1, LastModified: ModifiedOn, InsertedOn: Insertedtime},
	)
	if err != nil {
		fmt.Println("NOT ABLE TO ADD TO THE MONGODB", err)
	}
}

//UPDATE SCHEDULES
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

//UPDATE
func Update(session *mgo.Session, taskJs map[string]interface{}, lastmodified time.Time) {
	session.SetMode(mgo.Monotonic, true)

	Col := session.DB("TskSch").C("Schedule")

	id := int(taskJs["id"].(float64))
	name := taskJs["name"].(string)
	task := taskJs["cmd"].(string)
	week := int(taskJs["week"].(float64))
	day := int(taskJs["day"].(float64))
	second := int(taskJs["second"].(float64))
	minute := int(taskJs["minute"].(float64))
	hour := int(taskJs["hour"].(float64))
	r := int(taskJs["r"].(float64))

	SelectFrom := bson.M{"id": id}

	ChangeTo := bson.M{"$set": bson.M{"name": name, "task": task, "hour": hour, "minute": minute, "second": second, "day": day, "week": week, "r": r, "update": 1, "lastmodified": lastmodified}}

	Err := Col.Update(SelectFrom, ChangeTo)
	if Err != nil {
		fmt.Println("NOT ABLE TO UPDATE TO THE MONGODB", Err)
	}
}

//UPDATE
func UUpdate(session *mgo.Session, id int, task string, hour int, minute int, second int, day int, week int, r int, updated int, lastmodified time.Time) {
	session.SetMode(mgo.Monotonic, true)

	Col := session.DB("TskSch").C("Schedule")

	SelectFrom := bson.M{"id": id}

	ChangeTo := bson.M{"$set": bson.M{"task": task, "hour": hour, "minute": minute, "second": second, "day": day, "week": week, "r": r, "update": updated, "lastmodified": lastmodified}}

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
