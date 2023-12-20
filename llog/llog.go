package llog

import (
	"io"
	"log"
	"os"
	"wols/config"
)

func levelLog(level int, logString string) {
	if config.Cfg.EnableLog {
		logFile, err := os.OpenFile(config.Cfg.LogFile, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
		if err != nil {
			panic(err)
		}
		log.SetOutput(io.MultiWriter(os.Stderr, logFile))
		defer logFile.Close()
	}

	l := []string{"DEBUG", "INFO", "WARN", "ERROR", "PANIC", "FETAL"}
	log.Printf("[%v] %v\n", l[level-1], logString)
}

func Debug(String string) {
	if config.Cfg.LogLevel <= config.LvlDebug {
		levelLog(config.LvlDebug, String)
	}
}

func Info(String string) {
	if config.Cfg.LogLevel <= config.LvlInfo {
		levelLog(config.LvlInfo, String)
	}
}

func Warn(String string) {
	if config.Cfg.LogLevel <= config.LvlWarn {
		levelLog(config.LvlWarn, String)
	}
}

func Error(String string) {
	if config.Cfg.LogLevel <= config.LvlError {
		levelLog(config.LvlError, String)
	}
}

func Panic(String string) {
	if config.Cfg.LogLevel <= config.LvlPanic {
		levelLog(config.LvlPanic, String)
	}
}

func Fetal(String string) {
	if config.Cfg.LogLevel <= config.LvlFetal {
		levelLog(config.LvlFetal, String)
	}
}
