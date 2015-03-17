package schedule

import (
	"fmt"
	"gopkg.in/tomb.v2"
	"github.com/garyburd/redigo/redis"
	"gopkg.in/mgo.v2"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Schedule struct {
	L  string
	Id int
	W  *sync.WaitGroup
	Session *mgo.Session
	Conn redis.Conn
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
		Time := schedule[1]
		Cmd := schedule[2]
		t, _ := strconv.Atoi(Time)
		SchTime := time.Duration(t)
		for _ = range time.Tick(SchTime * time.Second) {
			func() {
				Sch.T.M.Lock()
				put2msgQ(Cmd,Sch.Conn,Sch.Session,Sch.Id)
				Sch.T.M.Unlock()
				fmt.Println("TASK", Sch.Id ,"GOT EXECUTED")
			}()
		}
	Sch.W.Done()
	return nil
}

func put2msgQ(Cmd string,Conn redis.Conn,session *mgo.Session,cmd_id int){

	//PUSHING THE cmd_id msgQ
	_, err := Conn.Do("LPUSH", "task", cmd_id)
	if err != nil {
		fmt.Println("CAN'T PUSH IT BACK")
	}else{
		put2resDB(session,cmd_id)
		put2Cmdtxt(Cmd,cmd_id) 
	}
}

func put2Cmdtxt(Cmd string, cmd_id int) {

	//INSERTING INTO command.txt
	file, err := os.OpenFile("../command.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("ERROR IN OPENING command.txt FILE")
	}
	line := strconv.Itoa(cmd_id) + ":" + Cmd + "\n"
	if _, Err := file.WriteString(line); Err != nil {
		fmt.Println("ERROR IN WRITING INTO command.txt FILE")
	}
}

func put2resDB(session *mgo.Session,cmd_id int){

	//INSERTING INTO RESULTDB
	session.SetMode(mgo.Monotonic, true)
	Col := session.DB("TskSch").C("Result")
	err := Col.Insert(&Result{strconv.Itoa(cmd_id),false,"","",0,false,"",""})
	if err !=nil {
		fmt.Println("NOT ABLE TO INSERT TO resultDB" , err)
	}

}

func (Sch *Schedule) Stop() error {
	Sch.T.Kill(nil)
	return Sch.T.Wait()
}
