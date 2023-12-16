package webs

import (
	"embed"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"wols/cmds"
	"wols/llog"
	"wols/nic"
	"wols/recent"
	"wols/wol"
)

type responseStates struct {
	Status string `json:"status"`
	Extra  string `json:"extra"`
}

func (rs *responseStates) json() []byte {
	bytes, _ := json.MarshalIndent(rs, "", "  ")
	return bytes
}

//go:embed static
var webFS embed.FS
var emb = true

func setMime(s string) (string, string) {
	mime := map[string]string{
		"html": "text/html",
		"htm":  "text/html",
		"css":  "text/css",
		"js":   "text/javascript",
		"ico":  "image/png",
		"json": "application/json",
	}

	s = filepath.Ext(s)
	if s == "" {
		s = ".html"
	}

	s = strings.ReplaceAll(s, ".", "")

	ts, ok := mime[s]
	if !ok {
		return "Content-Type", "application/octet-stream"
	}
	return "Content-Type", ts
}

func respHtml(w http.ResponseWriter, r *http.Request) {
	llog.Debug(r.URL.String())
	switch r.URL.Path {
	case "/":
		fallthrough
	case "/index.htm":
		fallthrough
	case "/index.html":
		w.Header().Set(setMime(r.URL.Path))

		b, err := webFS.ReadFile("static/index.html")
		if !emb {
			b, err = os.ReadFile("webs/static/index.html")
		}

		if err != nil {
			llog.Error(err.Error() + ":" + r.URL.Path)
		}

		fmt.Fprint(w, string(b))

	default:
		w.WriteHeader(http.StatusNotFound)
		llog.Debug("StatusNotFound(404) -> " + r.URL.Path)
	}
}

func broadCast(w http.ResponseWriter, r *http.Request) {
	var rs responseStates

	err := r.ParseForm()
	if err != nil {
		llog.Error("broadcast: " + err.Error())
	}
	mac := ""
	desc := ""

	for k, v := range r.Form {
		if k == "mac" {
			mac = strings.Join(v, "")
		}
		if k == "desc" {
			desc = strings.Join(v, "")
		}
	}
	llog.Debug(fmt.Sprintf("MAC: %v;DSC: %v", mac, desc))
	hwAddr, err := nic.StringToMAC(mac)
	if err != nil {
		llog.Error(err.Error())
		rs.Status = "error"
		rs.Extra = err.Error()
		w.Write(rs.json())
		return
	}

	wol.BroadcastMagicPack(hwAddr)
	_, err = recent.Add(hwAddr, desc)
	if err != nil {
		llog.Error("Error to add recent: " + err.Error())
		rs.Status = "error"
		rs.Extra = err.Error()
		w.Write(rs.json())
		return
	}

	rs.Status = "success"
	w.Write(rs.json())
}

func reMove(w http.ResponseWriter, r *http.Request) {
	var rs responseStates

	err := r.ParseForm()
	if err != nil {
		rs.Status = "error"
		rs.Extra = err.Error()
		w.Write(rs.json())
		return
	}
	mac := ""
	for k, v := range r.Form {
		if k == "mac" {
			mac = strings.Join(v, "")
		}
	}
	hwAddr, err := nic.StringToMAC(mac)
	if err != nil {
		llog.Error(err.Error())
		rs.Status = "error"
		rs.Extra = err.Error()
		w.Write(rs.json())
		return
	}
	err = recent.Remove(hwAddr)
	if err != nil {
		llog.Error("Error to remove recent: " + err.Error())
		rs.Status = "error"
		rs.Extra = err.Error()
		w.Write(rs.json())
		return
	}
	rs.Status = "success"
	w.Write(rs.json())
}

func moDify(w http.ResponseWriter, r *http.Request) {
	var rs responseStates

	err := r.ParseForm()
	if err != nil {
		llog.Error("modify: " + err.Error())
	}
	desc := ""
	mac := ""
	for k, v := range r.Form {
		if k == "mac" {
			mac = strings.Join(v, "")
		}

		if k == "desc" {
			desc = strings.Join(v, "")
		}
	}
	llog.Debug("modify desc: " + mac + desc)
	hwAddr, err := nic.StringToMAC(mac)
	if err != nil {
		llog.Error(err.Error())
		rs.Status = "error"
		rs.Extra = err.Error()
		w.Write(rs.json())
		return
	}
	_, err = recent.Modify(hwAddr, desc)
	if err != nil {
		llog.Error("modify desc:" + mac + " error:" + err.Error())
		rs.Status = "error"
		rs.Extra = err.Error()
		w.Write(rs.json())
		return
	}
	// todo: 返回信息而不是JSON
	rs.Status = "success"
	w.Write(rs.json())
}

func putJson(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(setMime(r.URL.Path))

	w.Write(recent.Json())
}

func respStatic(w http.ResponseWriter, r *http.Request) {
	s := r.URL.Path
	if s == "/favicon.ico" {
		s = "/favicon.png"
	}
	s = "static" + s

	b, err := webFS.ReadFile(s)
	if !emb {
		b, err = os.ReadFile("webs/" + s)
	}

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		llog.Error(err.Error() + ":" + r.URL.Path)
		return
	}
	w.Header().Set(setMime(s))
	w.Write(b)
}

func WEBServ() {
	llog.Info(fmt.Sprint("WEB Server listen on port:", strconv.Itoa(cmds.PortWebs)))

	http.HandleFunc("/", respHtml)
	http.HandleFunc("/broadcast", broadCast)
	http.HandleFunc("/remove", reMove)
	http.HandleFunc("/modify", moDify)
	http.HandleFunc("/recents", putJson)
	http.HandleFunc("/favicon.ico", respStatic)
	http.HandleFunc("/css/", respStatic)
	http.HandleFunc("/script/", respStatic)

	err := http.ListenAndServe(":"+strconv.Itoa(cmds.PortWebs), nil)
	if err != nil {
		llog.Error("ListenAndServe: " + err.Error())
	}

}
