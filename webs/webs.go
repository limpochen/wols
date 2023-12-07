package webs

import (
	"embed"
	"fmt"
	"html"
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

//go:embed static
var webFS embed.FS

var emb = false

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
	//ParseRequest(r)
	err := r.ParseForm() // 解析参数，默认是不会解析的
	if err != nil {
		println(err.Error())
	}
	mac := ""

	for k, v := range r.Form {
		if k == "mac" {
			mac = html.EscapeString(strings.Join(v, ""))
		}
	}

	hwAddr, err := nic.StringToMAC(mac)
	if err != nil {
		llog.Error(err.Error())
		return
	}

	wol.BroadcastMagicPack(hwAddr)
	recent.Add(hwAddr, "from Web")

	w.Header().Set(setMime(r.URL.Path))
	w.Write(recent.Json())
}

func reMove(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm() // 解析参数，默认是不会解析的
	if err != nil {
		println(err.Error())
	}
	mac := ""
	for k, v := range r.Form {
		if k == "mac" {
			mac = html.EscapeString(strings.Join(v, ""))
		}
	}
	hwAddr, err := nic.StringToMAC(mac)
	if err != nil {
		llog.Error(err.Error())
		return
	}
	recent.Remove(hwAddr)
	w.Header().Set(setMime(r.URL.Path))
	w.Write(recent.Json())
}

func moDify(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm() // 解析参数，默认是不会解析的
	if err != nil {
		println(err.Error())
	}
	desc := ""
	mac := ""
	for k, v := range r.Form {
		if k == "mac" {
			mac = html.EscapeString(strings.Join(v, ""))
		}

		if k == "desc" {
			desc = html.EscapeString(strings.Join(v, ""))
		}
	}
	hwAddr, err := nic.StringToMAC(mac)
	if err != nil {
		llog.Error(err.Error())
		return
	}
	_, err = recent.Modify(hwAddr, desc)
	if err != nil {
		llog.Error("modify desc:" + mac + " error:" + err.Error())
	}
	w.Header().Set(setMime(r.URL.Path))
	w.Write([]byte(desc))
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
	http.HandleFunc("/broadcast.html", broadCast)
	http.HandleFunc("/remove.html", reMove)
	http.HandleFunc("/modify.html", moDify)
	http.HandleFunc("/recents.json", putJson)
	http.HandleFunc("/favicon.ico", respStatic)
	http.HandleFunc("/css/", respStatic)
	http.HandleFunc("/purecss3/", respStatic)
	http.HandleFunc("/script/", respStatic)

	err := http.ListenAndServe(":"+strconv.Itoa(cmds.PortWebs), nil)
	if err != nil {
		llog.Error("ListenAndServe: " + err.Error())
	}

}

func ParseRequest(r *http.Request) {
	fmt.Printf("Mothod:\t%v\n", r.Method)

	fmt.Printf("URL:\t%v\n", r.URL)
	fmt.Printf("\tScheme:\t%v\n", r.URL.Scheme)
	fmt.Printf("\tOpaque:\t%v\n", r.URL.Opaque)
	fmt.Printf("\tUser:\t%v\n", r.URL.User)
	fmt.Printf("\tHost:\t%v\n", r.URL.Host)
	fmt.Printf("\tPath:\t%v\n", r.URL.Path)
	fmt.Printf("\tRawPath:\t%v\n", r.URL.RawPath)
	fmt.Printf("\tOmitHost:\t%v\n", r.URL.OmitHost)
	fmt.Printf("\tForceQuery:\t%v\n", r.URL.ForceQuery)
	fmt.Printf("\tRawQuery:\t%v\n", r.URL.RawQuery)
	fmt.Printf("\tFragment:\t%v\n", r.URL.Fragment)
	fmt.Printf("\tRawFragment:\t%v\n", r.URL.RawFragment)

	fmt.Println("Header:")
	for k, v := range r.Header {
		fmt.Printf("\t%v:\n", k)
		for _, v2 := range v {
			fmt.Printf("\t\t%v\n", v2)
		}
	}

	fmt.Printf("Host:\t%v\n", r.Host)

	r.ParseForm()
	fmt.Println("Form:")
	for k, v := range r.Form {
		fmt.Printf("\t%v\t%v:\n", k, v)
	}
	fmt.Println("PostForm:")
	for k, v := range r.PostForm {
		fmt.Printf("\t%v\t%v:\n", k, v)
	}
	fmt.Printf("RemoteAddr:\t%v\n", r.RemoteAddr)
	fmt.Printf("RequestURI:\t%v\n", r.RequestURI)

}
