package main
import (
        "github.com/gorilla/mux"
        "net/http"
        "strings"
        "strconv"
        "fmt"
        "os"
        "bytes"
        "io/ioutil"
        "mime/multipart"
        "TskSch/resultDB"
        "archive/zip"
        "io"
)
type a struct {
        Host string
        Port string
}
var agentInfo map[string]a
func main(){
    agentInfo = make(map[string]a)
    m := mux.NewRouter()

    //PING
    m.HandleFunc("/ping",func(w http.ResponseWriter, req *http.Request) {
        result := req.FormValue("name")
        if result == "manager"{
            w.WriteHeader(200)
            w.Write([]byte("{\"status\":\"alive\"}"))
        }else if result == "all"{
            var r []string
            for k, _ := range agentInfo{
                path := "http://"+agentInfo[k].Host+":"+agentInfo[k].Port+"/ping"
                res , err := http.Get(path)
                if err != nil {
                        fmt.Println("CAN'T CONNECT TO YOUR TASK AGENT : " + k )
                }
                body , _ := ioutil.ReadAll(res.Body)
                if string(body) != ""{
                        r = append(r,"\"" + k + "\": " + "\"alive\"")
                }else{
                        r = append(r,"\"" + k + "\": " + "\"notalive\"")
                }
            }
            str := strings.Join(r,",")
            w.WriteHeader(200)
            w.Write([]byte("{"+str+"}"))
        }else{
            agent := agentInfo[result]
            if agent.Host != "" && agent.Port != "" {
                    path := "http://"+agent.Host+":"+agent.Port+"/ping"
                    res , err := http.Get(path)
                    if err != nil {
                            fmt.Println("CAN'T CONNECT TO YOUR TASK AGENT")
                    }
                    body , _ := ioutil.ReadAll(res.Body)
                    if string(body) != ""{
                            w.WriteHeader(200)
                            w.Write([]byte("{\"status\":\"alive\"}"))
                    }else {
                            w.Write([]byte("{\"status\":\"notalive\"}"))
                    }
            }else{
                    w.Write([]byte("{\"status\":\"notvalid\"}"))
            }
        }
    }).Methods("GET")

    //TASK INFO
    m.HandleFunc("/tasks",func(w http.ResponseWriter ,req * http.Request){
        result := req.FormValue("owner")
        if result != "all" {
                var r []string
                for k ,_ := range agentInfo{
                	path := "http://"+agentInfo[k].Host+":"+agentInfo[k].Port+"/tasks"
                    res , err := http.Get(path)
                    if err != nil {
                            fmt.Println("CAN'T CONNECT TO YOUR TASK AGENT : " + k )
                    }
                    body , _ := ioutil.ReadAll(res.Body)
                    if string(body) != ""{
                            r = append(r,"\"" + k + "\": " + "\""+string(body)+"\"")
                    }else{
                            r = append(r,"\"" + k + "\": " + "\""+""+"\"")
                    }
                }
                str := strings.Join(r,",")
                w.WriteHeader(200)
                w.Write([]byte("{"+str+"}"))
        }
    }).Methods("GET")

    //REGISTER
    m.HandleFunc("/register",func(w http.ResponseWriter, req *http.Request){
        result := req.FormValue("agent")
        if result != "" {
                        agentData := a{
                                Host : strings.Split(result,":")[0],
                                Port : strings.Split(result,":")[1],
                        }
                        agentName := strings.Split(result,":")[2]
                        agentInfo[agentName] = agentData
                        w.Write([]byte("ok"))
        }else {
                 w.Write([]byte("notok"))
        }
    }).Methods("GET")

	//SEND TASK_CMD when Task Agent asks the Scheduler
    m.HandleFunc("/askCommand", func(w http.ResponseWriter, req *http.Request){

		result := req.FormValue("cmdId")
		name := req.FormValue("agentName")
		cmd_id := strings.Split(result,":")[0]
		target_url := agentInfo[name].Host+":"+agentInfo[name].Port
		if target_url != "" {
			if cmd_id != ""{
				val , _ := strconv.Atoi(cmd_id)
				cmd := resultDB.Find(val,strings.Split(result,":")[1])
				path := "/home/solution/go/src/TskSch/Schclient/"+strings.Split(cmd,":")[1]+"/"
				files, _ := ioutil.ReadDir(path)		
				zipFile, err := os.Create("/home/solution/go/src/TskSch/Schclient/"+strings.Split(cmd,":")[1]+".zip")
				if err != nil {
					fmt.Println(err)
				}
				r := zip.NewWriter(zipFile)
				for _ , f := range files {
					fds ,err := os.Open(path+f.Name())
					if err != nil {
						fmt.Println(err)
					}
					// Create a new zip archive.
					fd, err := r.Create(f.Name())
					if err != nil {
						fmt.Println(err)
					}
					io.Copy(fd, fds)
				}
				// Make sure to check the error on Close.
				err = r.Close()
				if err != nil {
					fmt.Println(err)
				}
				flag := post("/home/solution/go/src/TskSch/Schclient/"+strings.Split(cmd,":")[1]+".zip",target_url)
				if flag != nil  {
					http.Error(w, "file can't be uploaded to taskagent", http.StatusBadRequest)
				}
				w.Write([]byte(cmd))
			}else{
					http.Error(w, "cmd_id cannot be empty", http.StatusBadRequest)
			}
		}else{
				http.Error(w, "Invalid task agent name", http.StatusBadRequest)
		}
	}).Methods("GET")

    //RUNNING THE SERVER AT PORT 8000
    err := http.ListenAndServe(":8000", m)
    if err != nil {
            fmt.Println("Error starting server on port.")
            fmt.Println(err)
    }
}
func post( file string,targetUrl string) error {
    bodyBuf := &bytes.Buffer{}
    bodyWriter := multipart.NewWriter(bodyBuf)

    // this step is very important
    fileWriter, err := bodyWriter.CreateFormFile("uploadfile", file)
    if err != nil {
        fmt.Println("error writing to buffer")
        return err
    }

    // open file handle
    fh, err := os.Open(file)
    if err != nil {
        fmt.Println("error opening file")
        return err
    }

    //iocopy
    _, err = io.Copy(fileWriter, fh)
    if err != nil {
        return err
    }

    contentType := bodyWriter.FormDataContentType()
    bodyWriter.Close()

    resp, err := http.Post("http://"+targetUrl+"/upload", contentType, bodyBuf)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    return nil
}
