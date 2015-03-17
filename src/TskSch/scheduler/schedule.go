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
	var count int = 0
	for {
		Cursor := SchCol.Find(nil)
		Cursor.Skip(count)
		iter := Cursor.Iter()
		for iter.Next(&res){
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
			count = count + 1
		}
		fmt.Println(SchMap)
		select{
		case <-Sch.T.Dying():
				Sch.W.Done()
		}
	}
	Sch.W.Wait()
}
func Restart(task_id int,task string,time int,day int){
	fmt.Println(SchMap)
	SchMap[task_id].Stop()
	Sch := new(schedule.Schedule)
	Wg.Add(1)
	Sch = &schedule.Schedule{
				L : strconv.Itoa(day) + ":" + strconv.Itoa(time) + ":" + task,
				Id : task_id,
				W : &Wg,
			}
	SchMap[task_id]=Sch
	Sch.T.Go(Sch.Push)
	fmt.Println(task_id,"GOT RESTARTED")
	Sch.W.Wait()
}
