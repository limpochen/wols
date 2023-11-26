package cmds

import (
	"flag"
	"fmt"
)

/*
Usage:	wols [-Flags] [-Flags] ...
Flags:

	-hwaddr <MAC>:	 	MAC to be broadcast.
	-cycle <int>:		Broadcast cycle (1 to 16).
	-no-wols:			Disable WOL service.
	-no-webs:			Disable WEB service.
	-no-scan:			Do not scan LAN hosts.
	-port-wols <int>:	Port of WOL service.
	-port-webs <int>:	Port of WEB service.
	-port-send <int>:	Port of Boradcast.

Eg:
	wols send -hwaddr 12:34:56:70:9A:BC
	wols serv -no-wols -port-webs 8080
*/

var HWAddr string
var BCCycle int
var NoWols, NoWebs, NoScan bool
var PortWols, PortWebs, PortSent int

func init() {
	flag.StringVar(&HWAddr, "hwaddr", "", "MAC to be broadcast.")
	flag.IntVar(&BCCycle, "cycle", 3, "Broadcast cycle (1 to 16).")
	flag.BoolVar(&NoWols, "no-wols", false, "Disable WOL service.")
	flag.BoolVar(&NoWebs, "no-webs", false, "Disable WEB service.")
	flag.BoolVar(&NoScan, "no-scan", false, "Do not scan local network hosts.")
	flag.IntVar(&PortWols, "port-wols", 12307, "Port of WOL service.")
	flag.IntVar(&PortWebs, "port-webs", 7077, "Port of WEB service.")
	flag.IntVar(&PortSent, "port-sent", 7, "Port of Boradcast.")

	flag.Parse()
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
