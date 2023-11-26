package wol

import (
	"fmt"
	"net"
	"strconv"
	"wols/cmds"
)

func WOLServ() {
	//设置UDP监听地址
	addr, err := net.ResolveUDPAddr("udp", "0.0.0.0:"+strconv.Itoa(cmds.PortWols))
	if err != nil {
		panic(err)
	}
	//开始UDP监听
	fmt.Println("WOL Server listen on port:" + strconv.Itoa(cmds.PortWols))
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	//接收UDP数据
	for {
		// Here must use make and give the lenth of buffer
		bufUDP := make([]byte, 60000)
		_, _, err := conn.ReadFromUDP(bufUDP)
		if err != nil {
			fmt.Println(err)
			continue
		}
		hwAddr := GetMagicPacketMacFromBuffer(bufUDP)
		if hwAddr != nil {
			BroadcastMagicPack(*hwAddr)
		}
	}
}
