package nic

import (
	"fmt"
	"math/rand"
	"net"
	"wols/config"
)

type Nip struct {
	IPv   int       `json:"-"`
	Ip    net.IP    `json:"ip"`
	IpNet net.IPNet `json:"-"`
	Arps  []Arp     `json:"arps,omitempty"`
}

type Nif struct {
	Index  int              `json:"index"`
	Name   string           `json:"name"`
	Mac    net.HardwareAddr `json:"-"`
	MacStr string           `json:"mac"`
	Nips   []Nip            `json:"nips,omitempty"`
}

var Nifs []Nif

func (n Nip) GetBroadcastIP() net.IP {
	if len(n.Ip) != 4 { //len(n.IpNet.Mask)
		return nil
	}

	broadcast := net.IP(make([]byte, len(n.Ip)))

	for idx := range n.Ip {
		broadcast[idx] = n.Ip[idx] | ^n.IpNet.Mask[idx]
	}
	return broadcast
}

func (n Nip) IsBroadcastIP(ip net.IP) bool {
	if len(ip) != 4 {
		ip = ip.To4()
	}
	if ip == nil {
		return false
	}

	ipCast := n.GetBroadcastIP()
	if ipCast == nil {
		return false
	}

	for i := range ipCast {
		if ipCast[i] != ip[i] {
			return false
		}
	}
	return true
}

func ParseNif() error {
	Nifs = Nifs[0:0]

	ns, err := net.Interfaces()
	if err != nil {
		return err
	}

	for i := range ns {
		var n Nif
		var nips []Nip

		if ns[i].Flags&net.FlagRunning != net.FlagRunning {
			continue
		}

		if ns[i].Flags&net.FlagLoopback == net.FlagLoopback {
			continue
		}

		n.Index = ns[i].Index
		n.Name = ns[i].Name
		n.Mac = ns[i].HardwareAddr
		n.MacStr = n.Mac.String()

		addrs, _ := ns[i].Addrs()
		for j := range addrs {
			var nip Nip

			ip, ipnet, _ := net.ParseCIDR(addrs[j].String())

			ip4 := ip.To4()
			if ip4 == nil {
				continue
			}
			nip.Ip = ip4

			nip.IPv = len(nip.Ip)
			nip.IpNet = *ipnet

			if nip.Ip.IsLoopback() || nip.Ip.IsUnspecified() {
				break
			}
			nips = append(nips, nip)
		}
		n.Nips = nips
		Nifs = append(Nifs, n)
	}

	return nil
}

func (n Nif) Print() {
	fmt.Printf("Index:%v\tName:%v\tMAC:%v\n", n.Index, n.Name, n.Mac)
	for i := range n.Nips {
		fmt.Printf("\tIP: %v\n", n.Nips[i].Ip)
		fmt.Printf("\t\tNetIP:\t%v\n", n.Nips[i].IpNet.IP)
		fmt.Printf("\t\tMask:\t%v\n", n.Nips[i].IpNet.Mask)
		fmt.Printf("\t\tCast:\t%v\n", n.Nips[i].GetBroadcastIP())
		fmt.Printf("\t\tArps:\n")
		for ii := range n.Nips[i].Arps {
			fmt.Printf("\t\t\t%v\t%v\n", n.Nips[i].Arps[ii].Ip.String(), n.Nips[i].Arps[ii].Mac.String())
		}
	}
}

func Test() []byte {
	ln := 1433
	st := 711
	mac, _ := StringToMAC(config.HWAddr)

	buf := make([]byte, ln)
	for i := 0; i < ln; i++ {
		buf[i] = byte(rand.Intn(255))
	}
	for i := 0; i < 6; i++ {
		buf[st+i] = 0xff
	}
	for i := 0; i < 16; i++ {
		for j := 0; j < 6; j++ {
			buf[st+i*6+j+6] = mac[j]
		}
	}
	return buf
}

/*
type Nic struct {
	IPv    int              `json:"-"`
	Index  int              `json:"INDEX"`
	Name   string           `json:"NAME"`
	Mac    net.HardwareAddr `json:"-"`
	MacStr string           `json:"MAC"`
	Ip     net.IP           `json:"IP"`
	IpNet  net.IPNet        `json:"-"`
	Arps   []Arps           `json:"ARPS,omitempty"`
}

var Nics []Nic

func (n Nic) GetBroadcastIP() net.IP {
	if len(n.Ip) != 4 { //len(n.IpNet.Mask)
		return nil
	}

	broadcast := net.IP(make([]byte, len(n.Ip)))

	for idx := range n.Ip {
		broadcast[idx] = n.Ip[idx] | ^n.IpNet.Mask[idx]
	}
	return broadcast
}

func (n Nic) IsBroadcastIP(ip net.IP) bool {
	if len(ip) != 4 {
		ip = ip.To4()
	}
	if ip == nil {
		return false
	}

	ipCast := n.GetBroadcastIP()
	if ipCast == nil {
		return false
	}

	for i := range ipCast {
		if ipCast[i] != ip[i] {
			return false
		}
	}
	return true
}

func ParseNic() error {
	var ni Nic
	Nics = Nics[0:0]

	ns, err := net.Interfaces()
	if err != nil {
		return err
	}

	for i := range ns {
		if ns[i].Flags&net.FlagRunning != net.FlagRunning {
			continue
		}
		ni.Index = ns[i].Index
		ni.Name = ns[i].Name
		ni.Mac = ns[i].HardwareAddr
		ni.MacStr = ni.Mac.String()
		addrs, _ := ns[i].Addrs()
		for j := range addrs {
			ip, ipnet, _ := net.ParseCIDR(addrs[j].String())

			ip4 := ip.To4()
			if ip4 != nil {
				ni.Ip = ip4
			} else {
				ni.Ip = ip
			}
			ni.IPv = len(ni.Ip)
			ni.IpNet = *ipnet

			if ni.Ip.IsLoopback() || ni.Ip.IsUnspecified() {
				break
			}
			Nics = append(Nics, ni)
		}
	}
	fmt.Printf("%v\n", Nics)
	return nil
}

func (n Nic) Print() {
	fmt.Printf(":======\nIndex:\t%v\n", n.Index)
	fmt.Printf("Name:\t%v\n", n.Name)
	fmt.Printf("IP:\t%v\n", n.Ip)
	fmt.Printf("MAC:\t%v\n", n.Mac)
	fmt.Printf("IPNet->IP:\t%v\n", n.IpNet.IP)
	fmt.Printf("IPNet->Mask:\t%v\n", n.IpNet.Mask)
	fmt.Printf("Cast:\t%v\n", n.GetBroadcastIP())
	fmt.Printf("Arps:\n")
	for i := range n.Arps {
		fmt.Printf("\t%v\t%v\n", n.Arps[i].Ip.String(), n.Arps[i].Mac.String())
	}
}
*/
