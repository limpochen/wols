package config

import (
	"flag"
	"fmt"
	"os"
)

var (
	HWAddr   string
	HomePath string
)

func init() {

	flag.StringVar(&HWAddr, "hwaddr", "", "MAC to be broadcast")
	flag.StringVar(&HomePath, "home", "", "Config file path")

	flag.Parse()
}

func Usage() error {
	fmt.Print("wols (Wake-On-Lan Integrated service tool) Copyright(C) 2023 limpo@live.com\n\n")
	if HomePath == "" {
		HomePath = ExecPath
	} else {
		if _, err := os.Stat(HomePath); os.IsNotExist(err) {
			err := os.MkdirAll(HomePath, 0755)
			if err != nil {
				fmt.Println("Error creating folder:", err)
				os.Exit(1)
			}
		}
	}
	return nil
}
