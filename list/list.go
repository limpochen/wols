package list

import (
	"time"
	"wols/nic"
)

func ScanLAN() {
	nic.ParseNif()
	nic.SystemArp()
	//for i := range nic.Nifs {
	//	nic.Nifs[i].Print()
	//}

	for {
		time.Sleep(time.Second * 10)
	}
}
