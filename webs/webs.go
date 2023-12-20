package webs

import (
	"embed"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"wols/config"
	"wols/llog"
	"wols/nic"
	"wols/recent"
	"wols/wol"
)

//go:embed static
var webFS embed.FS
var emb = false

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

	switch opts[2] {
	case "broadcast":
		wol.BroadcastMagicPack(hwAddr)
		_, err = recent.Add(hwAddr, strings.Join(r.Form["desc"], ""))
		if err != nil {
			llog.Error("Error to add recent: " + err.Error())
			rs.Status = "error"
			rs.Extra = err.Error()
			w.Write(rs.json())
			return
		}
	case "remove":
		err = recent.Remove(hwAddr)
		if err != nil {
			llog.Error("Error to remove recent: " + err.Error())
			rs.Status = "error"
			rs.Extra = err.Error()
			w.Write(rs.json())
			return
		}

	case "modify":
		_, err = recent.Modify(hwAddr, strings.Join(r.Form["desc"], ""))
		if err != nil {
			llog.Error("modify desc:" + err.Error())
			rs.Status = "error"
			rs.Extra = err.Error()
			w.Write(rs.json())
			return
		}
	default:
		llog.Error("unknown command")
		rs.Status = "error"
		rs.Extra = "unknown command"
		w.Write(rs.json())
		return
	}

	rs.Status = "success"
	w.Write(rs.json())
}

func WEBServ() {
	llog.Info("WEB Server listen on port:" + strconv.Itoa(config.Cfg.WebsPort))

	http.HandleFunc("/opt/", basicAuth(respOpt))
	http.HandleFunc("/", basicAuth(respStatic))
	http.HandleFunc("/recents", basicAuth(putJson))

	err := http.ListenAndServe(":"+strconv.Itoa(config.Cfg.WebsPort), nil)
	if err != nil {
		llog.Error("ListenAndServe: " + err.Error())
	}

}

func basicAuth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !config.RequireAuth() {
			next.ServeHTTP(w, r)
			return
		}
		username, password, ok := r.BasicAuth()
		if ok && config.AuthUser(username, password) {
			next.ServeHTTP(w, r)
			return
		}

		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
}
