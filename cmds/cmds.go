package cmds

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var (
	ExecPath                     string
	BaseName                     string
	LogFile                      string
	LogLevelString               string
	NoLog                        bool
	HWAddr                       string
	BCCycle                      int
	NoWols, NoWebs, NoScan       bool
	PortWols, PortWebs, PortSent int
)
var LogLevel = -1 //Undefined

func init() {
	ExecPath, _ = os.Executable()
	ExecPath, _ = filepath.EvalSymlinks(ExecPath)
	ext := filepath.Ext(ExecPath)
	BaseName = strings.TrimSuffix(ExecPath, ext)
	ExecPath = filepath.Dir(ExecPath)

	flag.StringVar(&HWAddr, "hwaddr", "", "MAC to be broadcast.")
	flag.IntVar(&BCCycle, "cycle", 3, "Broadcast cycle (1 to 16).")
	flag.BoolVar(&NoWols, "no-wols", false, "Disable WOL service.")
	flag.BoolVar(&NoWebs, "no-webs", false, "Disable WEB service.")
	flag.BoolVar(&NoScan, "no-scan", false, "Do not scan local network hosts.")
	flag.IntVar(&PortWols, "port-wols", 12307, "Port of WOL service.")
	flag.IntVar(&PortWebs, "port-webs", 7077, "Port of WEB service.")
	flag.IntVar(&PortSent, "port-sent", 7, "Port of Boradcast.")
	flag.StringVar(&LogFile, "logfile", "", "Full path of log file.")
	flag.BoolVar(&NoLog, "nolog", false, "Do not log to file.")
	flag.StringVar(&LogLevelString, "loglevel", "info", "debug, info, warn, error")
	flag.Parse()

	if LogFile == "" {
		LogFile = BaseName + ".log"
	}

	l := []string{"debug", "info", "warn", "error", "panic", "fetal"}
	for idx, lvl := range l {
		if LogLevelString == lvl {
			LogLevel = idx
			break
		}
	}
	if LogLevel == -1 {
		LogLevel = 1 //Info
	}
	println(LogLevel)
}

func Usage() error {

	if PortWebs == PortWols {
		return fmt.Errorf("the port of the WOL and WEB service cannot be the same")
	}

	if PortWols == 7 || PortWols == 9 {
		return fmt.Errorf("the WOL service listening port is incorrectly set. It must be not 7 and 9")
	}

	if PortWebs == 7 || PortWebs == 9 {
		return fmt.Errorf("the WEB service listening port is incorrectly set. It must be not 7 and 9")
	}

	if PortSent != 7 && PortSent != 9 {
		return fmt.Errorf("the broadcast port is incorrectly set. It must be 7 or 9")
	}

	if BCCycle < 1 || BCCycle > 16 {
		return fmt.Errorf("magic packet broadcast cycle incorrectly(1 to 16)")
	}

	if NoWols && NoWebs {
		return fmt.Errorf("the WOL or WEB service needs to enable at least one")
	}

	return nil
}
