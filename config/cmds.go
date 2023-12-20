package config

import (
	"flag"
	"fmt"
)

var (
	HWAddr string
)

func init() {

	flag.StringVar(&HWAddr, "hwaddr", "", "MAC to be broadcast")
	flag.StringVar(&ConfigFile, "c", "", "Config file path")

	flag.Parse()
}

func Usage() error {
	fmt.Print("wols (Wake-On-Lan Integrated service tool)  Copyright(C) 2023  limpo@live.com\n\n")

	return nil
}
