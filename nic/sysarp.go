package nic

import (
	"encoding/hex"
	"fmt"
	"net"
	"os/exec"
	"runtime"
	"strings"
)

type Arp struct {
	Ip     net.IP            `json:"ip"`
	Mac    HardwareAddrFixed `json:"-"`
	MacStr string            `json:"mac"`
}

func SystemArp() (arps []Arp) {
	switch runtime.GOOS {
	case "linux":
		arps = linuxArp()

	case "windows":
		arps = windowsArp()
	//case "darwin":

	default:
		return nil
	}

	return arps
}

func windowsArp() (arps []Arp) {
	//c := exec.Command("cmd", "/c", "chcp", "65001")
	//c.Run()
	c := exec.Command("arp", "-a")
	out, err := c.Output()
	if err != nil {
		fmt.Println(err)
		return nil
	}

	lines := string(out)

	idx := 0
	for _, v := range strings.Split(lines, "\n") {
		v = strings.Replace(v, "\r", "", -1)
		a := strings.Fields(v)

		if len(a) == 4 {
			if a[2] == "---" {
				idx = IndexOfNIC(a[3])
			}
		}

		if len(a) != 3 {
			continue
		}

		ip := net.ParseIP(a[0])
		if ip == nil {
			continue
		}

		mac, err := StringToMAC(a[1])
		if err != nil {
			continue
		}

		for i := range Nifs {
			for ii := range Nifs[i].Nips {
				if Nifs[i].Index == idx {
					if Nifs[i].Nips[ii].IPv != 4 {
						continue
					}
					if !Nifs[i].Nips[ii].IpNet.Contains(ip) {
						continue
					}
					if Nifs[i].Nips[ii].IsBroadcastIP(ip) {
						continue
					}
					Nifs[i].Nips[ii].Arps = append(Nifs[i].Nips[ii].Arps, Arp{ip, mac, mac.String()})
				}
			}
		}
	}
	return arps
}

func linuxArp() (arps []Arp) {
	cmd := exec.Command("arp", "-e")
	out, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	lines := string(out)
	for _, v := range strings.Split(lines, "\n") {
		v = strings.Replace(v, "\r", "", -1)
		a := strings.Fields(v)

		if len(a) != 5 {
			continue
		}

		ip := net.ParseIP(a[0])
		if ip == nil {
			continue
		}

		mac, err := StringToMAC(a[2])
		if err != nil {
			continue
		}

		for i := range Nifs {
			for ii := range Nifs[i].Nips {
				if Nifs[i].Name == a[4] {
					if Nifs[i].Nips[ii].IPv != 4 {
						continue
					}
					if !Nifs[i].Nips[ii].IpNet.Contains(ip) {
						continue
					}
					if Nifs[i].Nips[ii].IsBroadcastIP(ip) {
						continue
					}
					Nifs[i].Nips[ii].Arps = append(Nifs[i].Nips[ii].Arps, Arp{ip, mac, mac.String()})
				}
			}
		}
	}
	return arps
}

func IndexOfNIC(a string) int {
	a = strings.ToUpper(a)
	a = strings.ReplaceAll(a, "0X", "")
	c := 0
	if len(a) == 1 {
		a = "0" + a
	}
	if len(a) > 2 {
		return c
	}
	b, err := hex.DecodeString(a)
	if err != nil {
		return c
	}
	for _, v := range b {
		c = int(v)
	}

	return c
}
