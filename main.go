//go:generate goversioninfo
package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
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
		fmt.Fprintf(os.Stderr, "Error loading configuration: %v\n", err)
		return
	}

	if config.HWAddr != "" {
		hwAddr, err := nic.StringToMAC(config.HWAddr)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		wol.BroadcastMagicPack(hwAddr, "Command line")
		return
	}

	if !config.Cfg.Wols.EnableWols || !config.Cfg.Webs.EnableWebs {
		fmt.Fprintln(os.Stderr, "No services are enabled, Modify the configuration file to enable it.")
		return
	}

	if err = llog.LevelLog(config.LvlInfo, "WOLS started."); err != nil {
		fmt.Println("Log to file:", err)
		return
	}

	err = recent.Load()
	if err != nil {
		llog.Debug(err.Error())
	}

	config.HttpPort = config.Cfg.Webs.WebsPort + 10
	config.HttpsPort = config.Cfg.Webs.WebsPort + 20

	chWols := make(chan string)
	chWebs := make(chan string)
	chProxy := make(chan string)

	if config.Cfg.Wols.EnableWols {
		go wol.WOLServ(chWols)
	}

	if config.Cfg.Webs.EnableWebs {
		go webs.WEBServ(chWebs)
	}

	//if !config.NoScan {
	//	go list.ScanLAN()
	//}

	c = make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	for {
		interrupt := false
		select {
		case status := <-chWols:
			if status == "ok" {
				llog.Info("WOL Server listen on port:" + strconv.Itoa(config.Cfg.Wols.WolsPort))
			}
			if status == "error" {
				interrupt = true
			}
		case status := <-chWebs:
			if status == "ok" {
				llog.Info("WEB Server listen on port:" + strconv.Itoa(config.Cfg.Webs.WebsPort))
			}
			if status == "error" || status == "shutdown" {
				interrupt = true
			}
		case status := <-chProxy:
			if status == "error" {
				interrupt = true
			}

		case <-c:
			interrupt = true
		}

		if interrupt {
			llog.Warn("wols shutting down.")
			break
		}

	}

}
