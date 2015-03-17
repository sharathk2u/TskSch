package scheduler
import(
	"fmt"
	"sync"
	"TskSch/schedule"
	"TskSch/resultDB"
	"strconv"
	"time"
	"gopkg.in/mgo.v2"
	"github.com/garyburd/redigo/redis"
)
type SchResult struct {
	Id int
	Task string
	Time int
	Day int
	Updated int
	LastModified time.Time
	InsertedOn time.Time
}

var SchMap map[int]*schedule.Schedule
var Sch *schedule.Schedule
var Wg     sync.WaitGroup

func Schedule(Session *mgo.Session,conn redis.Conn){
	SchMap = make(map[int]*schedule.Schedule)
	Sch = new(schedule.Schedule)
	var res = &SchResult{}
	session := resultDB.ResultdbInit()
	session.SetMode(mgo.Monotonic, true)
	SchCol := session.DB("TskSch").C("Schedule")
	for {
		Cursor := SchCol.Find(nil)
		iter := Cursor.Iter()
		for iter.Next(&res){
			if(res.Updated == 1){
				if(res.LastModified.Equal(res.InsertedOn)){
				Wg.Add(1)
				Sch = &schedule.Schedule{
					L : strconv.Itoa(res.Day) + ":" + strconv.Itoa(res.Time) + ":" + res.Task,
					Id : res.Id,
					W : &Wg,
					Session : Session,
					Conn : conn,
				}
				Sch.T.Go(Sch.Push)
				SchMap[Sch.Id] = Sch
				resultDB.UpdateSchedule(session,res.Id,0)
				fmt.Println(res.Id,"STARTED")
				}else{
					resultDB.UpdateSchedule(session,res.Id,0)
					Restart(res.Id , res.Task , res.Time , res.Day)
				}
			}
		}
	}
	Sch.W.Wait()
}

func Restart(task_id int,task string,time int,day int){
	SchMap[task_id].T.Kill(fmt.Errorf(strconv.Itoa(task_id),"UPDATED"))
	Sch := new(schedule.Schedule)
	Wg.Add(1)
	Sch = &schedule.Schedule{
				L : strconv.Itoa(day) + ":" + strconv.Itoa(time) + ":" + task,
				Id : task_id,
				W : &Wg,
			}
	SchMap[task_id]=Sch
	Sch.T.Go(Sch.Push)
	fmt.Println(task_id,"RESTARTED")
}
