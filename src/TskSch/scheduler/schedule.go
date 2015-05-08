package scheduler
import(
	"fmt"
	"sync"
	"TskSch/schedule"
	"TskSch/resultDB"
	"strconv"
	"time"
	"gopkg.in/mgo.v2"
)
type SchResult struct {
	Id int
	Name string
	Task string
	Hour int
	Minute int
	Second int
	Day int
	Week int
	R int
	Update int
	LastModified time.Time
	InsertedOn time.Time
}

var SchMap map[int]*schedule.Schedule
var Sch *schedule.Schedule
var Wg     sync.WaitGroup

func Schedule(Session *mgo.Session,host1 string ,host string ,port string){
	SchMap = make(map[int]*schedule.Schedule)
	Sch = new(schedule.Schedule)
	var res = &SchResult{}
	session := resultDB.ResultdbInit(host1)
	session.SetMode(mgo.Monotonic, true)
	SchCol := session.DB("TskSch").C("Schedule")
	for {
		Cursor := SchCol.Find(nil)
		iter := Cursor.Iter()
		for iter.Next(&res){
			if(res.Update == 1){
				if(res.LastModified.Equal(res.InsertedOn)){
				Wg.Add(1)
				Sch = &schedule.Schedule{
					L : strconv.Itoa(res.R) + ":" +strconv.Itoa(res.Week) + ":" + strconv.Itoa(res.Day) + ":" +strconv.Itoa(res.Hour) + ":" + strconv.Itoa(res.Minute) +":"+ strconv.Itoa(res.Second) + ":" + res.Task,
					Id : res.Id,
					Name : res.Name,
					W : &Wg,
					Session : Session,
					Host : host,
					Port : port,
				}
				Sch.T.Go(Sch.Push)
				SchMap[Sch.Id] = Sch
				//fmt.Println(SchMap)
				resultDB.UpdateSchedule(session,res.Id,0)
				fmt.Println(res.Id,"STARTED")
				}else{
					resultDB.UpdateSchedule(session,res.Id,0)
					//fmt.Println(SchMap)
					Restart(Session ,res.Id, res.Name, res.Task , res.R ,res.Week, res.Day , res.Hour , res.Minute , res.Second)
				}
			}
		}
	}
	Sch.W.Wait()
}

func Restart(Session *mgo.Session ,task_id int,name string,task string, r int , week int, day int , hour int, minute int, second int){
	SchMap[task_id].T.Kill(fmt.Errorf(strconv.Itoa(task_id),"UPDATED"))
	Sch := new(schedule.Schedule)
	Wg.Add(1)
	Sch = &schedule.Schedule{
				L : strconv.Itoa(r) + ":" +strconv.Itoa(week) + ":" + strconv.Itoa(day) + ":" +strconv.Itoa(hour) + ":" + strconv.Itoa(minute) +":"+ strconv.Itoa(second) + ":" + task,
				Id : task_id,
				Name : name,
				W : &Wg,
				Session : Session,
			}
	SchMap[task_id]=Sch
	Sch.T.Go(Sch.Push)
	fmt.Println(task_id,"RESTARTED")
}
