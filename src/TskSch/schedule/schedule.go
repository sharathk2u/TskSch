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
	R , _ := strconv.Atoi(schedule[0])
	Week , _ := strconv.Atoi(schedule[1])
	Day , _ := strconv.Atoi(schedule[2])
	Hour, _ := strconv.Atoi(schedule[3])
	Minute , _ := strconv.Atoi(schedule[4])
	Second , _ := strconv.Atoi(schedule[5])
	Cmd := schedule[6]
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
		put2resDB(session,cmd_id)
//		put2Cmdtxt(Cmd,cmd_id) 
	}
}

//func put2Cmdtxt(Cmd string, cmd_id int) {

//	//INSERTING INTO command.txt
//	file, err := os.OpenFile("../command.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
//	if err != nil {
//		fmt.Println("ERROR IN OPENING command.txt FILE")
//	}
//	line := strconv.Itoa(cmd_id) + ":" + Cmd + "\n"
//	if _, Err := file.WriteString(line); Err != nil {
//		fmt.Println("ERROR IN WRITING INTO command.txt FILE")
//	}
//}

func put2resDB(session *mgo.Session,cmd_id int){

	//INSERTING INTO RESULTDB
	session.SetMode(mgo.Monotonic, true)
	Col := session.DB("TskSch").C("Result")
	err := Col.Insert(&Result{strconv.Itoa(cmd_id),false,"","",0,false,"",""})
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
