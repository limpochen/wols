package recent

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"
	"wols/config"
	"wols/nic"
)

type recent struct {
	Index    int                   `json:"index"`
	hwAddr   nic.HardwareAddrFixed `json:"-"`
	MAC      string                `json:"mac"`
	Count    int                   `json:"count"`
	lastTime time.Time             `json:"-"`
	Last     string                `json:"last"`
	Desc     string                `json:"desc,omitempty"`
}

var recents []recent

func Load() error {
	rf, err := os.ReadFile(config.Cfg.RecentsFile)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(rf, &recents); err != nil {
		return err
	}
	return Check()
}

func Check() error {
	for i, rs := range recents {
		hwAddr, err := nic.StringToMAC(rs.MAC)
		if err != nil {
			return fmt.Errorf("load recent error on index %v:%v", rs.Index, err)
		}
		recents[i].hwAddr = hwAddr

		last, err := time.Parse("2006-01-02 15:04:05", rs.Last)
		if err != nil {
			return fmt.Errorf("load recent last time error on index %v:%v", rs.Index, err)
		}
		recents[i].lastTime = last
	}
	return nil
}

func Write() error {
	if len(recents) == 0 {
		if err := os.Remove(config.Cfg.RecentsFile); err != nil {
			return err
		}
		return nil
	}

	if err := os.WriteFile(config.Cfg.RecentsFile, Json(), 0644); err != nil {
		return err
	}
	return nil
}

func Add(mac nic.HardwareAddrFixed, desc string) (index int, err error) {
	var r recent
	for i, r := range recents {
		if r.hwAddr == mac {
			recents[i].Count++
			recents[i].lastTime = time.Now()
			recents[i].Last = recents[i].lastTime.Format("2006-01-02 15:04:05")
			err := Write()
			if err != nil {
				return i, err
			}
			return i, nil
		}
	}
	r.Index = len(recents) + 1
	r.hwAddr = mac
	r.MAC = mac.String()
	r.Desc = desc
	r.Count = 1
	r.lastTime = time.Now()
	r.Last = r.lastTime.Format("2006-01-02 15:04:05")
	recents = append(recents, r)
	Cut256()
	err = Write()
	if err != nil {
		return 0, err
	}
	return 0, nil
}

func Modify(mac nic.HardwareAddrFixed, desc string) (index int, err error) {
	for i, r := range recents {
		if r.hwAddr == mac {
			recents[i].Desc = desc
			err = Write()
			if err != nil {
				return 0, err
			}
			return i, nil
		}
	}
	return len(recents), fmt.Errorf("recent modify desc error: MAC not matched")
}

func Remove(mac nic.HardwareAddrFixed) error {
	for i := range recents {
		if recents[i].hwAddr == mac {
			recents = append(recents[:i], recents[i+1:]...)

			err := Write()
			if err != nil {
				return err
			}
			return nil
		}
	}
	return fmt.Errorf("recent remove item error: MAC not matched")
}

func Cut256() error {
	if len(recents) < 256 {
		return nil
	}
	sort.Slice(recents, func(i, j int) bool {
		return recents[i].Last > recents[j].Last
	})
	toprecents := recents[:255]
	for i := range toprecents {
		toprecents[i].Index = i
	}
	recents = toprecents
	toprecents = nil
	err := Write()
	if err != nil {
		return err
	}
	return nil
}

func Json() []byte {
	bytes, _ := json.MarshalIndent(recents, "", "  ")

	return bytes
}
