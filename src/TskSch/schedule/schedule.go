package schedule

import (
	"fmt"
	"sync"
	"gopkg.in/mgo.v2"
	"github.com/garyburd/redigo/redis"
	"time"
	"strconv"
	"strings"
	"os"
	"os/exec"
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

func Push(Wg *sync.WaitGroup,line string,cmd_id int,ScheduleFile string,session *mgo.Session,Conn redis.Conn) {
	schedule := strings.Split(line,":")
	removeSchedule := "sed -i -e '1d' " + ScheduleFile
	Err := exec.Command("sh", "-c",removeSchedule).Run()
	if Err != nil {
		fmt.Println("DID NOT DELETE THE SCHEDULE FROM schedule.txt")
	}
	if len(schedule) == 2 {
		Time := schedule[0]
		Cmd := schedule[1]
		t , _ := strconv.Atoi(Time)
		SchTime := time.Duration(t)
		for _ = range time.Tick(SchTime * time.Second) {
			func(){
				fmt.Println("DOING UR ACTIONS...")
				put2msgQ(Cmd,Conn,session,cmd_id)
				fmt.Println("ACTIONS COMPLETE !!!")
			}()
		}
	}
	if len(schedule) == 3 {
		Day := schedule[0]
		Time := schedule[1]
		Cmd := schedule[2]
		d , _ := strconv.Atoi(Day)
		t , _ := strconv.Atoi(Time)
		t = t + (d * 24 * 60 * 60)
		SchTime := time.Duration(t)
		for _ = range time.Tick(SchTime * time.Second) {
			func(){
				fmt.Println("DOING UR ACTIONS...")
				put2msgQ(Cmd,Conn,session,cmd_id)
				fmt.Println("ACTIONS COMPLETE !!!")
			}()
		}
	}
	Wg.Done()
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

func put2Cmdtxt(Cmd string,cmd_id int ){

	//INSERTING INTO command.txt
	file, err := os.OpenFile("../command.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Error opening command.txt file", err)
		os.Exit(1)
	}
	line := strconv.Itoa(cmd_id) + ":" + Cmd + "\n"
	if _, err = file.WriteString(line); err != nil {
		fmt.Println(" CAN'T WRITE INTO FILE ")
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

