//go:generate goversioninfo
package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
	"wols/cmds"
	"wols/list"
	"wols/nic"
	"wols/webs"
	"wols/wol"
)

var c chan os.Signal

func main() {
	fmt.Print("wols (Wake-On-Lan Integrated service tool)  Copyright(C) 2023  limpo@live.com\n\n")
	err := cmds.Usage()
	if err != nil {
		fmt.Println(err)
		return
	}

	if cmds.HWAddr != "" {
		hwAddr, err := nic.StringToMAC(cmds.HWAddr)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		wol.BroadcastMagicPack(hwAddr)
		return
	}

	logFile, err := os.OpenFile(cmds.LogFile, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}
	log.SetOutput(io.MultiWriter(os.Stderr, logFile))
	defer logFile.Close()

	if !cmds.NoWebs {
		go webs.WEBServ()
	}

	if !cmds.NoWols {
		go wol.WOLServ()
	}

	if !cmds.NoScan {
		go list.ScanLAN()
	}

	c = make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-c
}
