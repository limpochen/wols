package wol

import (
	"fmt"
	"net"
	"wols/config"
	"wols/llog"
	"wols/nic"
)

const MagicPacketLen = 102

// 唤醒魔术包，含6字节头部，16次重复 MAC 网络地址
type magicPacket struct {
	Header  [6]byte
	Payload [16]nic.HardwareAddrFixed
}

// 函数：生成魔术包
// 参数：网卡MAC地址
// 返回：魔术包结构体
func genMagicPacket(mac nic.HardwareAddrFixed) (packet magicPacket) {

	// 生成6个字节头部
	for idx := range packet.Header {
		packet.Header[idx] = 0xFF
	}

	// 填充16遍MAC地址
	for idx := range packet.Payload {
		packet.Payload[idx] = mac
	}

	return packet
}

func BroadcastMagicPack(hwAddr nic.HardwareAddrFixed) {
	nic.ParseNif()
	for i := range nic.Nifs {
		for ii := range nic.Nifs[i].Nips {
			if nic.Nifs[i].Nips[ii].IPv != 4 {
				continue
			}

			la := net.UDPAddr{
				IP: nic.Nifs[i].Nips[ii].Ip,
			}
			ra := net.UDPAddr{
				IP:   nic.Nifs[i].Nips[ii].GetBroadcastIP(),
				Port: config.Cfg.BroadcastPort,
			}
			c, err := net.DialUDP("udp", &la, &ra)
			if err != nil {
				llog.Error(fmt.Sprint(err, ":", nic.Nifs[i].Nips[ii].Ip))
				//return err
				continue
			}

			for idx := 1; idx <= config.Cfg.BroadcastCycle; idx++ {
				_, err = c.Write(genMagicPacket(hwAddr).Bytes())
				if err != nil {
					llog.Error(err.Error())
				}
			}
			c.Close()
			llog.Info(fmt.Sprintf("from %v broadcast %v at %s:%d",
				nic.Nifs[i].Nips[ii].Ip.String(),
				hwAddr.String(),
				nic.Nifs[i].Nips[ii].GetBroadcastIP().String(),
				config.Cfg.BroadcastPort))
		}
	}
}

// 从魔术包中提取硬件MAC地址
// 参数：魔术包结构体
// 返回：硬件MAC地址
func (packet magicPacket) GetMac() (mac nic.HardwareAddrFixed, err error) {
	for idx := range packet.Header {
		if packet.Header[idx] != 0xFF {
			return mac, fmt.Errorf("bad MagicPacket header")
		}
	}
	for idx := range packet.Payload {
		if idx == 0 {
			mac = packet.Payload[idx]
		} else {
			if packet.Payload[idx] != mac {
				return mac, fmt.Errorf("bad MagicPacket")
			}
		}
	}
	return mac, nil
}

func (packet magicPacket) Bytes() []byte {
	buf := make([]byte, MagicPacketLen)
	for i := range buf {
		if i < 6 {
			buf[i] = 0xFF
		} else {
			buf[i] = packet.Payload[0][i%nic.HwAddrFixedLen]
		}
	}
	return buf
}

func GetMagicPacketMacFromBuffer(buf []byte) *nic.HardwareAddrFixed {
	var mac nic.HardwareAddrFixed

	if len(buf) < MagicPacketLen {
		return nil
	}

	for idx := 0; idx < len(buf)-101; idx++ {
		checkOK := false
		PayloadWrong := false
		mac_all_f := 0
		mac_all_z := 0
		// 校验头部
		for i := 0; i < 6; i++ {
			if buf[idx+i] != 0xff {
				break
			}
			if i == 5 {
				checkOK = true
			}
		}
		if !checkOK {
			continue
		}

		// 头部命中，检验并提取16次MAC的第1次
		for i := 0; i < nic.HwAddrFixedLen; i++ {
			mac[i] = buf[idx+i+6]
			switch buf[idx+i+6] {
			case 0xff:
				mac_all_f++
			case 0:
				mac_all_z++
			}
		}
		// 判断取得的 MAC 是否全 0xFF 或全 0
		if mac_all_f == nic.HwAddrFixedLen || mac_all_z == nic.HwAddrFixedLen {
			continue
		}

		// MAC 校验
		checkOK = false
		for i := 0; i < nic.HwAddrFixedLen; i++ {
			// 校验15次
			for j := 1; j < 16; j++ {
				if buf[idx+i+6] != buf[idx+i+j*nic.HwAddrFixedLen+6] {
					PayloadWrong = true
					break
				}
			}
			// 主体校验失败，继续下次循环
			if PayloadWrong {
				break
			}

			if i == nic.HwAddrFixedLen-1 {
				checkOK = true
			}
		}
		// 主体校验失败，继续下次循环
		if PayloadWrong {
			continue
		}

		if checkOK {
			return &mac
		}
	}

	return nil
}

func GetMagicpacket(buf []byte) *nic.HardwareAddrFixed {
	return nil
}
