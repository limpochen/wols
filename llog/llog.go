package llog

import (
	"io"
	"log"
	"os"
	"wols/config"
)

func LevelLog(level int, logString string) error {
	if config.Cfg.Llog.EnableLog {
		logFile, err := os.OpenFile(config.Cfg.Llog.LogFile, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
		if err != nil {
			return err
		}
		log.SetOutput(io.MultiWriter(os.Stderr, logFile))
		defer logFile.Close()
	}

	l := []string{"DEBUG", "INFO", "WARN", "ERROR", "PANIC", "FETAL"}
	log.Printf("[%v] %v\n", l[level-1], logString)
	return nil
}

func Debug(String string) {
	if config.Cfg.Llog.LogLevel <= config.LvlDebug {
		LevelLog(config.LvlDebug, String)
	}
}

func Info(String string) {
	if config.Cfg.Llog.LogLevel <= config.LvlInfo {
		LevelLog(config.LvlInfo, String)
	}
}

func Warn(String string) {
	if config.Cfg.Llog.LogLevel <= config.LvlWarn {
		LevelLog(config.LvlWarn, String)
	}
}

func Error(String string) {
	if config.Cfg.Llog.LogLevel <= config.LvlError {
		LevelLog(config.LvlError, String)
	}
}

func Panic(String string) {
	if config.Cfg.Llog.LogLevel <= config.LvlPanic {
		LevelLog(config.LvlPanic, String)
	}
}

func Fetal(String string) {
	if config.Cfg.Llog.LogLevel <= config.LvlFetal {
		LevelLog(config.LvlFetal, String)
	}
}
