package cmds

import (
	"fmt"
	"log"
	"os"
)

const (
	Info = iota
	Warning
	Error
)

var MyLogFile string

// Define a simple log processing method
// Param : logType uint32, logMsg string
// return: error
func MyLog(logType uint32, logMsg string) error {
	logFile, err := os.OpenFile(MyLogFile, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err == nil {
		//log.SetOutput(io.MultiWriter(os.Stderr, logFile))
		log.SetOutput(logFile)
		defer logFile.Close()
	}

	switch logType {
	case Info:
		log.Printf("[INF] %s\n", logMsg)
	case Warning:
		log.Printf("[WRN] %s\n", logMsg)
	case Error:
		log.Printf("[ERR] %s\n", logMsg)
	default:
		err = fmt.Errorf("error log type")
	}

	return err
}
