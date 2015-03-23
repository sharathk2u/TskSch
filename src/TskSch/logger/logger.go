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

//INITIALIZER FOR LOG FILE
func LogInit() *os.File {

	file, err := os.OpenFile("/home/unbxd/unbxd/src/TskSch/log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Error opening log.txt file", err)
		os.Exit(1)
	}
	return file
}
