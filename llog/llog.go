package llog

import (
	"log"
	"wols/cmds"
)

const (
	lDebug = iota
	lInfo
	lWarn
	lError
	lPanic
	lFetal
)

func levelLog(level int, logString string) {
	l := []string{"DEBUG", "INFO", "WARN", "ERROR", "PANIC", "FETAL"}
	log.Printf("[%v] %v\n", l[level], logString)
}

func Debug(String string) {
	if cmds.LogLevel <= lDebug {
		levelLog(lDebug, String)
	}
}

func Info(String string) {
	if cmds.LogLevel <= lInfo {
		levelLog(lInfo, String)
	}
}

func Warn(String string) {
	if cmds.LogLevel <= lWarn {
		levelLog(lWarn, String)
	}
}

func Error(String string) {
	if cmds.LogLevel <= lError {
		levelLog(lError, String)
	}
}

func Panic(String string) {
	if cmds.LogLevel <= lPanic {
		levelLog(lPanic, String)
	}
}

func Fetal(String string) {
	if cmds.LogLevel <= lFetal {
		levelLog(lFetal, String)
	}
}
