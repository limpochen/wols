package webs

import (
	"embed"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"wols/llog"
	"wols/nic"
	"wols/recent"
	"wols/wol"
)

//go:embed static
var webFS embed.FS
var emb = true

type responseStates struct {
	Status string `json:"status"`
	Extra  string `json:"extra"`
}

func (rs *responseStates) json() []byte {
	bytes, _ := json.MarshalIndent(rs, "", "  ")
	return bytes
}

func setMime(s string) (string, string) {
	mime := map[string]string{
		"html":   "text/html",
		"htm":    "text/html",
		"css":    "text/css",
		"js":     "text/javascript",
		"ico":    "image/png",
		"png":    "image/png",
		"json":   "application/json",
		"recent": "application/json",
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

func putJson(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(setMime(r.URL.Path))
	w.Write(recent.Json())
}

func respStatic(w http.ResponseWriter, r *http.Request) {
	s := r.URL.Path
	switch s {
	case "/":
		fallthrough
	case "/index.htm":
		s = "/index.html"

	case "/favicon.ico":
		s = "/favicon.png"

	default:
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

func respOpt(w http.ResponseWriter, r *http.Request) {
	var rs responseStates

	err := r.ParseForm()
	if err != nil {
		rs.Status = "error"
		rs.Extra = err.Error()
		w.Write(rs.json())
		return
	}

	hwAddr, err := nic.StringToMAC(strings.Join(r.Form["mac"], ""))
	if err != nil {
		llog.Error(err.Error())
		rs.Status = "error"
		rs.Extra = err.Error()
		w.Write(rs.json())
		return
	}

	var opts = strings.Split(r.URL.Path, "/")

	desc := strings.Join(r.Form["desc"], "")

	switch opts[2] {
	case "broadcast":
		wol.BroadcastMagicPack(hwAddr, desc)
		_, err = recent.Add(hwAddr, desc)
		if err != nil {
			rs.Status = "error"
			rs.Extra = "Add recent: " + err.Error()
			w.Write(rs.json())
			llog.Error(rs.Extra)
			return
		}
	case "remove":
		err = recent.Remove(hwAddr)
		if err != nil {
			rs.Status = "error"
			rs.Extra = "Remove recent: " + err.Error()
			w.Write(rs.json())
			llog.Error(rs.Extra)
			return
		}

	case "modify":
		_, err = recent.Modify(hwAddr, desc)
		if err != nil {
			rs.Status = "error"
			rs.Extra = "Modify desc:" + err.Error()
			w.Write(rs.json())
			llog.Error(rs.Extra)
			return
		}
	default:
		rs.Status = "error"
		rs.Extra = "Unknown command: " + opts[2]
		w.Write(rs.json())
		llog.Error(rs.Extra)
		return
	}

	rs.Status = "success"
	w.Write(rs.json())
}
