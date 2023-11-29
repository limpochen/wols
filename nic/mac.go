package nic

import (
	"fmt"
	"regexp"
	"strconv"
)

const HwAddrFixedLen = 6

// MAC 网络地址 6 个字节
type HardwareAddrFixed [6]byte

func (mac HardwareAddrFixed) String() string {
	const macDigit = "0123456789ABCDEF"
	//const macDigit = "0123456789abcdef"

	if len(mac) == 0 {
		return ""
	}
	buf := make([]byte, 0, len(mac)*3-1)
	for i, b := range mac {
		if i > 0 {
			buf = append(buf, ':')
		}
		buf = append(buf, macDigit[b>>4])
		buf = append(buf, macDigit[b&0xF])
	}
	return string(buf)
}

// 函数：转换文本为MAC地址
// 参数：含有MAC地址的文本
// 返回：6字节MAC地址
func StringToMAC(strMAC string) (mac HardwareAddrFixed, err error) {
	delims := ":-"
	reMAC := regexp.MustCompile(`^([0-9a-fA-F]{2}[` + delims + `]){5}([0-9a-fA-F]{2})$`)

	if !reMAC.MatchString(strMAC) {
		return mac, fmt.Errorf("%s is not a IEEE 802 MAC-48 address", strMAC)
	}

	for i := 0; i < 6; i++ {
		res, _ := strconv.ParseInt(strMAC[i*3:i*3+2], 16, 16)
		mac[i] = byte(res)
	}

	return mac, nil
}
