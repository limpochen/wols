package config

import (
	_ "embed"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
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
	BroadcastPort  int    //default: 7
	BroadcastCycle int    //default: 3
	RecentsFile    string //default:
	Auth           struct {
		Authentication bool
		Username       string
		Password       string
	}
	Llog struct {
		EnableLog bool   //default: true
		LogFile   string //default:
		LogLevel  int    //default: info
	}
	Wols struct {
		EnableWols bool //default: trur
		WolsPort   int  //default: 12307
	}
	Webs struct {
		EnableWebs bool //default: true
		WebsPort   int  //default: 7077
		EnableTls  bool
		CertFile   string
		KeyFile    string
	}
}

var Cfg config

var (
	ExecPath   string
	BaseName   string
	ConfigFile string //default:
)

var (
	HttpPort  int
	HttpsPort int
)

func init() {
	ExecPath, _ = os.Executable()
	ExecPath, _ = filepath.EvalSymlinks(ExecPath)
	BaseName = filepath.Base(strings.TrimSuffix(ExecPath, filepath.Ext(ExecPath)))
	ExecPath = strings.ReplaceAll(filepath.Dir(ExecPath), "\\", "/")
}

func Load() error {
	var isChange = false

	ConfigFile = path.Join(HomePath, BaseName+".toml")

	_, err := os.Lstat(ConfigFile)
	if errors.Is(err, fs.ErrNotExist) {
		if err := os.WriteFile(ConfigFile, sampleToml, 0644); err != nil {
			return fmt.Errorf("create default profile %v", err)
		}
	}

	md, err := toml.DecodeFile(ConfigFile, &Cfg)
	if err != nil {
		return err
	}

	if !md.IsDefined("BroadcastPort") {
		Cfg.BroadcastPort = 7
	}
	if !md.IsDefined("BroadcastCycle") {
		Cfg.BroadcastCycle = 3
	}
	if !md.IsDefined("RecentsFile") || Cfg.RecentsFile == "" {
		Cfg.RecentsFile = path.Join(HomePath, BaseName+".recent")
	}

	if !md.IsDefined("Auth", "Authentication") {
		println("Authentication")
		Cfg.Auth.Authentication = false
	}
	if Cfg.Auth.Authentication {
		if Cfg.Auth.Username != "" && Cfg.Auth.Password != "" {
			if !isBcryptHash(Cfg.Auth.Password) {
				Cfg.Auth.Password = encodePassword(Cfg.Auth.Password)
				isChange = true
			}
		} else {
			return errors.New("username and password should be configured")
		}
	}

	if !md.IsDefined("Llog", "EnableLog") {
		Cfg.Llog.EnableLog = true
	}
	if !md.IsDefined("Llog", "LogFile") || Cfg.Llog.LogFile == "" {
		Cfg.Llog.LogFile = path.Join(HomePath, BaseName+".log")
	}
	if !md.IsDefined("Llog", "LogLevel") {
		Cfg.Llog.LogLevel = LvlInfo
	}

	if !md.IsDefined("Wols", "EnableWols") {
		Cfg.Wols.EnableWols = true
	}
	if !md.IsDefined("Wols", "WolsPort") {
		Cfg.Wols.WolsPort = 12307
	}

	if !md.IsDefined("Webs", "EnableWebs") {
		Cfg.Webs.EnableWebs = true
	}
	if !md.IsDefined("Webs", "WebsPort") {
		Cfg.Webs.WebsPort = 7077
	}

	if !md.IsDefined("Webs", "EnableTls") {
		Cfg.Webs.EnableTls = false
	}
	if Cfg.Webs.EnableTls {
		if Cfg.Webs.CertFile == "" || Cfg.Webs.KeyFile == "" {
			return errors.New("please specify the certificate and key file")
		}
		_, err := os.Lstat(Cfg.Webs.CertFile)
		if errors.Is(err, fs.ErrNotExist) {
			return fmt.Errorf("file not exist: %v", Cfg.Webs.CertFile)
		}
		_, err = os.Lstat(Cfg.Webs.KeyFile)
		if errors.Is(err, fs.ErrNotExist) {
			return fmt.Errorf("file not exist: %v", Cfg.Webs.KeyFile)
		}
	}

	if isChange {
		cf, err := os.OpenFile(ConfigFile, os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		toml.NewEncoder(cf).EncodeWithComments(Cfg, md)
		cf.Close()
	}

	return nil
}
