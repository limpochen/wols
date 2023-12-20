//go:generate goversioninfo
package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"wols/config"
	"wols/llog"
	"wols/nic"
	"wols/recent"
	"wols/webs"
	"wols/wol"
)

var c chan os.Signal

func main() {
	//fmt.Print("wols (Wake-On-Lan Integrated service tool)  Copyright(C) 2023  limpo@live.com\n\n")
	err := config.Usage()
	if err != nil {
		fmt.Println(err)
		return
	}

	if err = config.Load(); err != nil {
		llog.Error(err.Error())
	}

	if config.HWAddr != "" {
		hwAddr, err := nic.StringToMAC(config.HWAddr)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		wol.BroadcastMagicPack(hwAddr)
		return
	}

	err = recent.Load()
	if err != nil {
		llog.Debug(err.Error())
	}

	if config.Cfg.EnableWebs {
		go webs.WEBServ()
	}

	if config.Cfg.EnableWols {
		go wol.WOLServ()
	}

	//if !config.NoScan {
	//	go list.ScanLAN()
	//}

	c = make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-c
}
