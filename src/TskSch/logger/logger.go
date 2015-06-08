package logger

import (
	"fmt"
	"log"
	"os"
)

//INITIALIZING LOG FOR SUCCESS
func Success(file *os.File) *log.Logger {
	LogSucc := log.New(file, "SUCCESS: ", log.Ldate|log.Ltime|log.Lshortfile)
	return LogSucc
}

//INITIALIZING LOG FOR FAILURE
func Failure(file *os.File) *log.Logger {
	LogFail := log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	return LogFail
}

//INITIALIZING LOG FOR FAILURE
func Info(file *os.File) *log.Logger {
	LogInfo := log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	return LogInfo
}

//INITIALIZER FOR LOG FILE
func LogInit() *os.File {

	file, err := os.OpenFile("../log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Error opening log.txt file", err)
		os.Exit(1)
	}
	return file
}

//INITIALIZER FOR SCHEDULAR LOG FILE
func LogSchInit() *os.File {

	file, err := os.OpenFile("../log_schedular.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Error opening log.txt file", err)
		os.Exit(1)
	}
	return file
}

//INITIALIZER FOR MANAGER LOG FILE
func LogManInit() *os.File {

	file, err := os.OpenFile("../log_manager.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Error opening log.txt file", err)
		os.Exit(1)
	}
	return file
}

//INITIALIZER FOR MANAGER LOG FILE
func LogAgentInit() *os.File {

	file, err := os.OpenFile("../log_agent.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Error opening log.txt file", err)
		os.Exit(1)
	}
	return file
}