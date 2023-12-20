package config

import (
	_ "embed"
	"os"
	"path/filepath"
	"strings"

	"github.com/GuanceCloud/toml"
)

//go:embed sample.toml
var sampleToml []byte

const (
	LvlDebug = iota + 1
	LvlInfo
	LvlWarn
	LvlError
	LvlPanic
	LvlFetal
)

type config struct {
	Authentication bool
	Username       string
	Password       string
	EnableLog      bool   //default: true
	LogFile        string //default:
	LogLevel       int    //default: info
	RecentsFile    string //default:
	EnableWols     bool   //default: trur
	EnableWebs     bool   //default: true
	WolsPort       int    //default: 12307
	WebsPort       int    //default: 7077
	BroadcastPort  int    //default: 7
	BroadcastCycle int    //default: 3
}

var Cfg config

var (
	ExecPath   string
	BaseName   string
	ConfigFile string //default:
)

func init() {
	ExecPath, _ = os.Executable()
	ExecPath, _ = filepath.EvalSymlinks(ExecPath)
	ext := filepath.Ext(ExecPath)
	BaseName = strings.TrimSuffix(ExecPath, ext)
	ExecPath = filepath.Dir(ExecPath)

}

func Load() error {
	var configFile string
	var isChange = false

	if ConfigFile == "" {
		configFile = filepath.ToSlash(BaseName) + ".toml"
	}

	if !FileExist(configFile) {
		os.WriteFile(configFile, sampleToml, 0644)
	}

	md, err := toml.DecodeFile(configFile, &Cfg)
	if err != nil {
		return err
	}

	if !md.IsDefined("Authentication") {
		Cfg.Authentication = false
		isChange = true
	}

	if Cfg.Username != "" && Cfg.Password != "" && !isBcryptHash(Cfg.Password) {
		Cfg.Password = encodePassword(Cfg.Password)
		isChange = true
	}

	if !md.IsDefined("EnableLog") {
		Cfg.EnableLog = true
		isChange = true
	}

	if !md.IsDefined("LogFile") || Cfg.LogFile == "" {
		Cfg.LogFile = BaseName + ".log"
		isChange = true
	}

	if !md.IsDefined("LogLevel") {
		Cfg.LogLevel = LvlInfo
		isChange = true
	}

	if Cfg.RecentsFile == "" {
		Cfg.RecentsFile = BaseName + ".recent"
		isChange = true
	}

	if !md.IsDefined("EnableWols") {
		Cfg.EnableWols = true
		isChange = true
	}

	if !md.IsDefined("EnableWebs") {
		Cfg.EnableWebs = true
		isChange = true
	}

	if !md.IsDefined("WolsPort") {
		Cfg.WolsPort = 12307
		isChange = true
	}

	if !md.IsDefined("WebsPort") {
		Cfg.WebsPort = 7077
		isChange = true
	}

	if !md.IsDefined("BroadcastPort") {
		Cfg.BroadcastPort = 7
		isChange = true
	}

	if !md.IsDefined("BroadcastCycle") {
		Cfg.BroadcastCycle = 3
		isChange = true
	}

	if isChange {
		cf, err := os.OpenFile(configFile, os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			return err
		}
		toml.NewEncoder(cf).EncodeWithComments(Cfg, md)
		cf.Close()
	}

	return nil
}

func FileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err) //需要添加判断权限的操作
}
