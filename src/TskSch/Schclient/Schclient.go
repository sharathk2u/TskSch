package main
import(
    "TskSch/scheduler"
    "TskSch/resultDB"
    "encoding/json"
    "net/http"
    "fmt"
    "github.com/gorilla/mux"
    "time"
    "strconv"
    "strings"
    "code.google.com/p/goconf/conf"
	"TskSch/mailer"
	"os"
	"path/filepath"
//	"reflect"
//	"io"
	"bufio"
)

func main(){

        c, err := conf.ReadConfigFile("../server.conf")
        if err != nil {
                fmt.Println("CAN'T READ CONF FIILE",err)
        }
        host1 ,_ := c.GetString("resultDB","host")
        host2 ,_ := c.GetString("msgQ","host")
        port ,_ := c.GetString("msgQ","port")
        //INITIALIZING THE MONGODB
        Session := resultDB.ResultdbInit(host1)

        //CLOSING ALL THE CONNECTION
        defer func(){
                Session.Close()
        }()
        go scheduler.Schedule(Session,host1,host2,port)

        go listenServe(host1)

        select{}
}
func listenServe(host1 string) {

        m := mux.NewRouter()

        //INITIALIZING THE MONGODB
        Session := resultDB.ResultdbInit(host1)

        //CLOSING ALL THE CONNECTION
        defer func(){
                Session.Close()
        }()

        //PING
        m.HandleFunc("/ping",func(w http.ResponseWriter, req *http.Request) {

	        w.WriteHeader(200)
	        w.Write([]byte("{\"status\":\"alive\"}"))

        }).Methods("GET")
	
	//ADDING THE TASKS TO MONGODB
        m.HandleFunc("/addTask",func(w http.ResponseWriter, req *http.Request) {
			type TaskInfo struct{
                                Name string
                                Cmd string
                                Hour int
                                Minute int
                                Second int
                                Day int
                                Week int
                                R int
                        }
			var taskJs TaskInfo
			req.ParseMultipartForm(2000000)	
                        
			if len(req.MultipartForm.Value) != 0 && len(req.MultipartForm.Value)!=0 {
				hour, _ := strconv.Atoi(req.FormValue("hour"))
                                minute,_ := strconv.Atoi(req.FormValue("minute"))
                                second,_ := strconv.Atoi(req.FormValue("second"))
                                day,_ := strconv.Atoi(req.FormValue("day"))
                                week,_ := strconv.Atoi(req.FormValue("week"))
                                r,_ := strconv.Atoi(req.FormValue("r"))
				file := req.MultipartForm.File
				fileName := file["files"][0].Filename
				fileDir := req.Form["name"][0]
				err := os.Mkdir("." + string(filepath.Separator) + fileDir,0777)
        			if err != nil {
                			fmt.Println("Unable to create the directory for writing. Check your write access privilege",err)
      				}
				o, er := os.Create("." + string(filepath.Separator) + fileDir+string(filepath.Separator)+fileName)
        			if er != nil {
                			fmt.Println("Unable to create the file for writing. Check your write access privilege",er,o)
         			}
         			// write the content from POST to the file
        			fd , e := file["files"][0].Open()
         			if e != nil {
                			 fmt.Println(e)
         			}	
			 	r1 := bufio.NewReader(fd)
				s, e := Readln(r1)
				str := ""
		    		for e == nil {
					if s!=""{
		    				//fmt.Println(s)
						str += s+"\n"
					}
		    			s,e = Readln(r1)
		    		}
			//	_, err = io.Copy(o, fd)
         		//	if err != nil {
                	//		 fmt.Println(err)
         		//	}
				fmt.Println(str)
				o.Write([]byte(str))
				fmt.Println("File uploaded successfully : ")
				
				taskJs = TaskInfo{
                                        Name : req.Form["name"][0],
                                        Cmd : req.Form["cmd"][0],
                                        Hour : hour,
                                        Minute : minute,
                                        Second : second,
                                        Day : day,
                                        Week : week,
                                        R : r,
                                }
			//	js, _ := json.Marshal(taskJs)	
                        //        out := resultDB.InsertSchedule(Session,js)
                        //        w.Write([]byte( "{" + "\"status\"" + " : \"Inserted,\"" +"\"Id\""+" : "+ "\""+strconv.Itoa(out)+"\""+"}"))
                        //        mailer.Mail("GOSERVE: Regarding Task addition", taskJs.Name + " ADDED ")
                        }else{
                                http.Error(w, "taskData cannot be empty", http.StatusBadRequest)
                                mailer.Mail("GOSERVE: Regarding Task addition", "Unable to ADD " + taskJs.Name + " Please check the format of addition")
                        }
        }).Methods("POST")


	//UPDATE THE TASK
        m.HandleFunc("/updateTask", func(w http.ResponseWriter, req *http.Request){

			type TaskInfo struct{
				Id string
                                Name string
                                Cmd string
                                Hour int
                                Minute int
                                Second int
                                Day int
                                Week int
                                R int
                        }
                        var taskJs TaskInfo
                        req.ParseForm()
                        if len(req.Form) != 0 {
                                hour, _ := strconv.Atoi(req.Form["hour"][0])
                                minute,_ := strconv.Atoi(req.Form["minute"][0])
                                second,_ := strconv.Atoi(req.Form["second"][0])
                                day,_ := strconv.Atoi(req.Form["day"][0])
                                week,_ := strconv.Atoi(req.Form["week"][0])
                                r,_ := strconv.Atoi(req.Form["r"][0])
                                taskJs = TaskInfo{
					Id : req.Form["id"][0],
                                        Name : req.Form["name"][0],
                                        Cmd : req.Form["cmd"][0],
                                        Hour : hour,
                                        Minute : minute,
                                        Second : second,
                                        Day : day,
                                        Week : week,
                                        R : r,
                                }
                                js, _ := json.Marshal(taskJs)

				resultDB.Update(Session,js,time.Now())
				w.Write([]byte( "{" + "\"status\"" + " : \"updated\""+"}"))
				mailer.Mail("GOSERVE: Regarding Task updation", taskJs.Name + " UPDATED ")
			}else{
					http.Error(w, "taskData cannot be empty", http.StatusBadRequest)
					mailer.Mail("GOSERVE: Regarding Task updation", "Unable to UPDATE " + taskJs.Name + " Please check the format of addition ")
			}

        }).Methods("GET")

        //SEND TASK_CMD when Task Agent asks the Scheduler
        m.HandleFunc("/askCommand", func(w http.ResponseWriter, req *http.Request){

			result := req.FormValue("cmdId")
			cmd_id := strings.Split(result,":")[0]
			if cmd_id != ""{
					val , _ := strconv.Atoi(cmd_id)
					cmd := resultDB.Find(val,strings.Split(result,":")[1])
					w.Write([]byte(cmd))
			}else{
					http.Error(w, "cmd_id cannot be empty", http.StatusBadRequest)
			}

        }).Methods("GET")

        //RUNNING THE SERVER AT PORT 8001
        err := http.ListenAndServe(":8001", m)
        if err != nil {
            fmt.Println("Error starting server on port.",err)
        }
}
func Readln(r *bufio.Reader) (string, error) {
        var (
                isRead bool = true
                Err error = nil
                line, ln []byte
        )
        for isRead && Err == nil {
                line, isRead , Err = r.ReadLine()
                ln = append(ln, line...)
	}
        return string(ln),Err
}
