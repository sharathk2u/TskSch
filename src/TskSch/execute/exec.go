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
type res struct{
	Value bool
	Name string
} 
type taskInfo struct {
        taskId map[string]res
        mutex *sync.Mutex
}

var Info = func() *taskInfo{
        return &taskInfo{taskId : make(map[string]res) , mutex : new(sync.Mutex)}
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
        taskname := cmd[2]

        inTime := time.Now()

        //UPDATING THE RESULTDB
        resultDB.UpdateResult(session,task_id,taskname,false, inTime.String(), time.Since(inTime).String(), os.Getpid(), true, "", "")

        //STORING TASKID IN taskInfo
        Info.mutex.Lock()
        r := res{
        		Value : true,
        		Name : taskname,		
        }
        Info.taskId[task_id]=r
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
        Remov(Info,task_id,taskname)

        if Err != nil {
                fmt.Println("ERROR IN EXECUTING THE COMMAND")
        }
        if Errout.String() != "" {
                //UPDATING THE RESULTDB
                resultDB.UpdateResult(session,task_id,taskname ,false, toe.String(), tte.String(), os.Getpid(), false, out.String(), Errout.String())
                //DUMPING THE RESULT TO LOG FILE
                LogFail := logger.Failure(file)
                LogFail.Println(Errout.String())
                fmt.Println("\t\tERROR\n", Errout.String())
        } else {
                //UPDATING THE RESULTDB
                resultDB.UpdateResult(session, task_id,taskname ,true, toe.String(), tte.String(), os.Getpid(), false, out.String(), Errout.String())
                //DUMPING THE RESULT TO LOG FILE
                LogSucc := logger.Success(file)
                LogSucc.Println("COMMAND GOT EXECUTED")
                fmt.Println(out.String())
        }
        Wg.Done()
}

func Remov(r *taskInfo,t_id string,name string) {
        r.mutex.Lock()
        x := res{
        	Value :false,
        	Name : name,
        }
        r.taskId[t_id] = x
        r.mutex.Unlock()
}

func Get() map[string]res {
        return Info.taskId
}
