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
        "TskSch/logger"
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

    //INITIALIZING THE LOG FILE
    logfile := logger.LogManInit()
    
    //CLOSING ALL THE CONNECTION
    defer func(){
            logfile.Close()
    }()

    //PING
    m.HandleFunc("/ping",func(w http.ResponseWriter, req *http.Request) {
        result := req.FormValue("name")
        if result == "manager"{
            w.WriteHeader(200)
            w.Write([]byte("{\"status\":\"alive\"}"))
            LogInfo := logger.Info(logfile)
            LogInfo.Println("MANAGER GOT PINGED")
        }else if result == "all"{
            var r []string
            for k, _ := range agentInfo{
                path := "http://"+agentInfo[k].Host+":"+agentInfo[k].Port+"/ping"
                res , err := http.Get(path)
                if err != nil {
                    fmt.Println("CAN'T CONNECT TO YOUR TASK AGENT : " + k )
                    LogErr := logger.Failure(logfile)
                    LogErr.Println("CAN'T CONNECT TO YOUR TASK AGENT : " + k )
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
            LogInfo := logger.Info(logfile)
            LogInfo.Println("EACH TASK AGENT GOT PINGED")
        }else{
            agent := agentInfo[result]
            if agent.Host != "" && agent.Port != "" {
                    path := "http://"+agent.Host+":"+agent.Port+"/ping"
                    res , err := http.Get(path)
                    if err != nil {
                        fmt.Println("CAN'T CONNECT TO TASK AGENT: " + result)
                        LogErr := logger.Failure(logfile)
                        LogErr.Println("CAN'T CONNECT TO TASK AGENT: "+ result)
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
            LogInfo := logger.Info(logfile)
            LogInfo.Println("TASK AGENT : "+ result+" GOT PINGED")
        }
    }).Methods("GET")

    //TASK INFO
    m.HandleFunc("/tasks",func(w http.ResponseWriter ,req * http.Request){
	result := req.FormValue("owner")
	var r []string
	if result == "all" {
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
        }else{
		path := "http://"+agentInfo[result].Host+":"+agentInfo[result].Port+"/tasks"
                    res , err := http.Get(path)
                    if err != nil {
                            fmt.Println("CAN'T CONNECT TO YOUR TASK AGENT : " + result )
                    }
                    body , _ := ioutil.ReadAll(res.Body)
                    if string(body) != ""{
                            r = append(r,"\"" + result + "\": " + "\""+string(body)+"\"")
                    }else{
                            r = append(r,"\"" + result + "\": " + "\""+""+"\"")
                    }
	
                str := strings.Join(r,",")
                w.WriteHeader(200)
                w.Write([]byte("{"+str+"}"))
	}	
    }).Methods("GET")

    //REGISTER
    m.HandleFunc("/register",func(w http.ResponseWriter, req *http.Request){
        LogInfo := logger.Info(logfile)
        result := req.FormValue("agent")
        var agentName string
        if result != "" {
            agentData := a{
                Host : strings.Split(result,":")[0],
                Port : strings.Split(result,":")[1],
            }
            agentName = strings.Split(result,":")[2]
            agentInfo[agentName] = agentData
            w.Write([]byte("ok"))
            LogInfo.Println("TASK AGENT :" + agentName+ " REGISTERED")
        }else {
            w.Write([]byte("notok"))
            LogInfo.Println("TASK AGENT :" + agentName+ " NOT REGISTERED")
        }
    }).Methods("GET")

	//SEND TASK_CMD when Task Agent asks the Scheduler
    m.HandleFunc("/askCommand", func(w http.ResponseWriter, req *http.Request){

		result := req.FormValue("cmdId")
		name := req.FormValue("agentName")
		cmd_id := strings.Split(result,":")[0]
		target_url := agentInfo[name].Host+":"+agentInfo[name].Port
		LogInfo := logger.Info(logfile)
        schedulerPath , _ := os.Getwd()
        if target_url != "" {
			if cmd_id != ""{
				val , _ := strconv.Atoi(cmd_id)
				cmd := resultDB.Find(val,strings.Split(result,":")[1])
				path := schedulerPath + strings.Split(cmd,":")[1]+"/"
				files, _ := ioutil.ReadDir(path)		
				zipFile, err := os.Create(schedulerPath + strings.Split(cmd,":")[1]+".zip")
				if err != nil {
					fmt.Println("Unable to create zipfile ",err)
				    LogErr := logger.Failure(logfile)
                    LogErr.Println("Unable to create zipfile ",err)
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
				LogInfo.Println("MANAGER CREATED THE " + strings.Split(cmd,":")[1] + ".zip" + " FILE")
                
                // Make sure to check the error on Close.
				err = r.Close()
				if err != nil {
					fmt.Println(err)
                    LogErr := logger.Failure(logfile)
                    LogErr.Println(err)
				}

                LogInfo.Println("MANAGER POSTING THE ZIPPED FILE TO THE TASK AGENT : " + target_url )
				flag := post(schedulerPath + strings.Split(cmd,":")[1]+".zip",target_url,logfile)
				if flag != nil  {
					http.Error(w, "file can't be uploaded to taskagent", http.StatusBadRequest)
                    LogErr := logger.Failure(logfile)
                    LogErr.Println("file can't be uploaded to taskagent")
				}else{
                    LogInfo.Println("ZIP FILE UPLOADED successfully TO TASK AGENT : " + target_url)    
					os.Remove(schedulerPath + strings.Split(cmd,":")[1]+".zip")
				    LogInfo.Println("ZIP FILE REMOVED successfully INSIDE MANAGER ")
                }
				w.Write([]byte(cmd))
                LogInfo.Println("MANAGER successfully GAVE THE COMMAND AND ZIP FILE TO TASK AGENT: " + target_url)
			}else{
				http.Error(w, "cmd_id cannot be empty", http.StatusBadRequest)
                LogErr := logger.Failure(logfile)
                LogErr.Println("cmd_id cannot be empty bad request from taskagent: "+ target_url )
			}
		}else{
			http.Error(w, "Invalid task agent name", http.StatusBadRequest)
            LogErr := logger.Failure(logfile)
            LogErr.Println("Invalid task agent name bad request")
        }
	}).Methods("GET")

    //RUNNING THE SERVER AT PORT 8000
    err := http.ListenAndServe(":8000", m)
    if err != nil {
        fmt.Println("Error starting server on port.",err)
        LogErr := logger.Failure(logfile)
        LogErr.Println("Error starting server on port.",err)
    }
}
func post( file string,targetUrl string,logfile *os.File) error {
    bodyBuf := &bytes.Buffer{}
    bodyWriter := multipart.NewWriter(bodyBuf)

    // this step is very important
    fileWriter, err := bodyWriter.CreateFormFile("uploadfile", file)
    if err != nil {
        LogErr := logger.Failure(logfile)
        LogErr.Println("error writing to buffer while posting zip file")
        return err
    }

    // open file handle
    fh, err := os.Open(file)
    if err != nil {
        LogErr := logger.Failure(logfile)
        LogErr.Println("error opening zip file: ",file,"while posting zip file")
        return err
    }

    //iocopy
    _, err = io.Copy(fileWriter, fh)
    if err != nil {
        LogErr := logger.Failure(logfile)
        LogErr.Println("error copying zip file: ",file,"while posting zip file")
        return err
    }

    contentType := bodyWriter.FormDataContentType()
    bodyWriter.Close()

    resp, err := http.Post("http://"+targetUrl+"/upload", contentType, bodyBuf)
    if err != nil {
        LogErr := logger.Failure(logfile)
        LogErr.Println("error uploading zip file: ",file,"while posting zip file")
        return err
    }
    defer resp.Body.Close()
    return nil
}
