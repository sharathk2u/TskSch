package main
import (
        "github.com/gorilla/mux"
        "net/http"
        "strings"
        "fmt"
        "io/ioutil"
//      "code.google.com/p/goconf/conf"
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

        //RUNNING THE SERVER AT PORT 8000
        err := http.ListenAndServe(":8000", m)
        if err != nil {
                fmt.Println("Error starting server on port.")
                fmt.Println(err)
        }
}

