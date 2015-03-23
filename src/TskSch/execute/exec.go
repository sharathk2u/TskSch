package execute

import (
	"TskSch/logger"
	"TskSch/resultDB"
	"bytes"
	"fmt"
	"gopkg.in/mgo.v2"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

type taskInfo struct {
	taskId map[string]bool
	mutex *sync.Mutex
}

var Info = func() *taskInfo{
	return &taskInfo{taskId : make(map[string]bool) , mutex : new(sync.Mutex)}
}()

func Exec(file *os.File, session *mgo.Session, Wg *sync.WaitGroup, args string) {

	//INITIALIZING STDOUT & STDERROR FOR COMMAND
	var (
		Errout bytes.Buffer
		out    bytes.Buffer
	)

	cmd := strings.Split(args, ":")
	task_id := cmd[0]
	task := cmd[1]

	inTime := time.Now()

	//UPDATING THE RESULTDB
	resultDB.UpdateResult(session, task_id, false, inTime.String(), time.Since(inTime).String(), os.Getpid(), true, "", "")

	//STORING TASKID IN taskInfo
	Info.mutex.Lock()
	Info.taskId[task_id]=true
	Info.mutex.Unlock()

	//EXECUTING THE COMMAND
	cmds := exec.Command("sh", "-c", task)
	cmds.Stderr = &Errout
	cmds.Stdout = &out

	//Time Of Execution => toe
	toe := time.Now()

	//Running the command
	Err := cmds.Run()

	//TIME TAKEN FOR EXECUTION
	tte := time.Since(toe)

	//REMOVING THE STORED INFO AFTER ITS EXECUTION IS COMPLETE
	Remov(Info,task_id)

	if Err != nil {
		fmt.Println("ERROR IN EXECUTING THE COMMAND")
	}
	if Errout.String() != "" {
		//UPDATING THE RESULTDB
		resultDB.UpdateResult(session, task_id, false, toe.String(), tte.String(), os.Getpid(), false, out.String(), Errout.String())
		//DUMPING THE RESULT TO LOG FILE
		LogFail := logger.Failure(file)
		LogFail.Println(Errout.String())
		fmt.Println("\t\tERROR\n", Errout.String())
	} else {
		//UPDATING THE RESULTDB
		resultDB.UpdateResult(session, task_id, true, toe.String(), tte.String(), os.Getpid(), false, out.String(), Errout.String())
		//DUMPING THE RESULT TO LOG FILE
		LogSucc := logger.Success(file)
		LogSucc.Println("COMMAND GOT EXECUTED")
		fmt.Println(out.String())
	}
	Wg.Done()
}

func Remov(r *taskInfo,t string) {
	r.mutex.Lock()
	r.taskId[t] = false
	r.mutex.Unlock()
}

func Get() map[string]bool {
	return Info.taskId
}
