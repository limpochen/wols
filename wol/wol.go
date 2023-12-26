package wol

import (
	"fmt"
	"net"
	"time"
	"wols/config"
	"wols/llog"
	"wols/nic"
)

func WOLServ(ch chan string) {
	//设置UDP监听地址
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", config.Cfg.Wols.WolsPort))
	if err != nil {
		llog.Error(err.Error())
		ch <- "error"
		return
	}
	//开始UDP监听
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		llog.Error(err.Error())
		ch <- "error"
		return
	}
	defer conn.Close()

	ch <- "ok"

	//接收UDP数据
	RCount := 0
	LastTime := time.Now()
	var LastMac nic.HardwareAddrFixed

	for {
		// Here must use make and give the lenth of buffer
		bufUDP := make([]byte, 60000)
		_, _, err := conn.ReadFromUDP(bufUDP)
		if err != nil {
			llog.Info(err.Error())
			continue
		}

		hwAddr := GetMagicPacketMacFromBuffer(bufUDP)
		if hwAddr == nil {
			continue
		}

		llog.Debug("recive MagicPacket: " + hwAddr.String())

		RCount++
		thisTime := time.Now()
		dur := thisTime.Sub(LastTime)
		LastTime = thisTime
		if dur < time.Duration(time.Millisecond*500) && LastMac == *hwAddr {
			llog.Debug("igrone the same magicpacket.")
			continue
		}

		BroadcastMagicPack(*hwAddr, "From WOLS")

		RCount = 0
		LastMac = *hwAddr
	}
}
