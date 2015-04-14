package schedule

import (
	"fmt"
	"gopkg.in/tomb.v2"
	"gopkg.in/mgo.v2"
	"strconv"
	"strings"
	"sync"
	"time"
	"TskSch/msgQ"
)

type Schedule struct {
	L  string
	Id int
	W  *sync.WaitGroup
	Session *mgo.Session
	T  tomb.Tomb
}

type Result struct {
	Task_id   string //command ID
	Cmd       string
	Executed  bool   //Executed staus
	TOE       string //Time Of Execution
	TTE       string //Time Taken to Execute
	Pid       int    // process id of client
	Exec_Stat bool   //Execution Status
	output    string
	err       string
}

func (Sch *Schedule) Push() error {
	schedule := strings.Split(Sch.L, ":")
	
	R , _ := strconv.Atoi(schedule[0])		// => For every interval or for only at particular time
	Week , _ := strconv.Atoi(schedule[1])	// => Week NAME
	Day , _ := strconv.Atoi(schedule[2])	// => For Every day : day = 1 or Foe Every 2nd day : day = 2
	Hour, _ := strconv.Atoi(schedule[3])	// => 24 Hr Format
	Minute , _ := strconv.Atoi(schedule[4]) // => Minutes
	Second , _ := strconv.Atoi(schedule[5])	// => Seconds
	
	Cmd := schedule[6]						// => Command

	if( R == 0 ){
		if(int(time.Now().Weekday()) == Week){
			ticker := updateTicker(Hour,Minute,Second,Day,7)
			for {
				<-ticker.C
				func() {
					Sch.T.M.Lock()
					put2msgQ(Cmd,Sch.Session,Sch.Id)
					Sch.T.M.Unlock()
					fmt.Println("TASK", Sch.Id ,"GOT EXECUTED")
				}()
				ticker = updateTicker(Hour,Minute,Second,Day,7)
			}
		}else{
			ticker := updateTicker(Hour,Minute,Second,Day,1)
			for {
				<-ticker.C
				func() {
					Sch.T.M.Lock()
					put2msgQ(Cmd,Sch.Session,Sch.Id)
					Sch.T.M.Unlock()
					fmt.Println("TASK", Sch.Id ,"GOT EXECUTED")
				}()
				ticker = updateTicker(Hour,Minute,Second,Day,1)
			}
		}
	}else {
		for _ = range time.Tick(time.Second*time.Duration( Hour*60 + Minute*60 + Second*1 )){
			if(Week == -1){
				func() {
					Sch.T.M.Lock()
					put2msgQ(Cmd,Sch.Session,Sch.Id)
					Sch.T.M.Unlock()
					fmt.Println("TASK", Sch.Id ,"GOT EXECUTED")
				}()
			}else{
				if(int(time.Now().Weekday()) == Week){
					func() {
						Sch.T.M.Lock()
						put2msgQ(Cmd,Sch.Session,Sch.Id)
						Sch.T.M.Unlock()
						fmt.Println("TASK", Sch.Id ,"GOT EXECUTED")
					}()
				}
			}
		}
	}
	Sch.W.Done()
	return nil
}

func put2msgQ(Cmd string,session *mgo.Session,cmd_id int){

	//INITIALIZING THE REDIS DB
	Conn := msgQ.RedisInit()
	defer func(){
		Conn.Close()
	}()
	//PUSHING THE cmd_id msgQ
	_, err := Conn.Do("LPUSH", "task", cmd_id)
	if err != nil {
		fmt.Println("CAN'T PUSH IT TO msgQ",err)
	}else{
		put2resDB(session,cmd_id,Cmd)
	}
}

func put2resDB(session *mgo.Session,cmd_id int , Cmd string){

	//INSERTING INTO RESULTDB
	session.SetMode(mgo.Monotonic, true)
	Col := session.DB("TskSch").C("Result")
	err := Col.Insert(&Result{strconv.Itoa(cmd_id),Cmd,false,"","",0,false,"",""})
	if err !=nil {
		fmt.Println("NOT ABLE TO INSERT TO resultDB" , err)
	}

}

func updateTicker(Hour int ,Minute int, Second int ,Day int,Week int ) *time.Ticker {
	nextTick := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), Hour, Minute, Second, 0, time.Local)
	if !nextTick.After(time.Now()) {
			nextTick = nextTick.Add(time.Duration(Week * Day * 24 ) * time.Hour)
	}
	diff := nextTick.Sub(time.Now())
	return time.NewTicker(diff)
}
