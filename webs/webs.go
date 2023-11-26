package webs

import (
	"embed"
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"wols/cmds"
	"wols/nic"
	"wols/wol"
)

//go:embed static
var webFS embed.FS

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
	switch r.URL.Path {
	case "/":
		fallthrough
	case "index.htm":
		fallthrough
	case "/index.html":
		w.Header().Set(setMime(r.URL.Path))
		r.ParseForm() // 解析参数，默认是不会解析的
		mac := ""
		msg := ""

		for k, v := range r.Form {
			if k == "mac" {
				mac = html.EscapeString(strings.Join(v, ""))
			}
		}

		if len(mac) != 0 {
			hwAddr, err := nic.StringToMAC(mac)
			if err != nil {
				fmt.Println(err)
				msg = fmt.Sprint(err)
			} else {
				wol.BroadcastMagicPack(hwAddr)
				msg = "MagicPacket sent to: " + mac
			}
		}

		b, err := webFS.ReadFile("static/index.html")
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERR: %v\n%v\n", err, "static"+r.URL.Path)
		}
		wf := string(b)
		wf = strings.Replace(wf, "@varmac@", mac, 1)
		wf = strings.Replace(wf, "@varmsg@", msg, 1)
		fmt.Fprint(w, wf)

	default:
		w.WriteHeader(http.StatusNotFound)
	}

}

func putJson(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(setMime(r.URL.Path))
	bytes, _ := json.MarshalIndent(nic.Nifs, "", "  ")
	//fmt.Print(string(bytes))
	w.Write(bytes)
}

func respStatic(w http.ResponseWriter, r *http.Request) {
	s := r.URL.Path
	if s == "/favicon.ico" {
		s = "/favicon.png"
	}
	b, err := webFS.ReadFile("static" + s)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(os.Stderr, "ERR: %v\n%v\n", err, "static"+r.URL.Path)
		return
	}
	w.Header().Set(setMime(s))
	w.Write(b)
}

func WEBServ() {
	fmt.Println("WEB Server listen on port:" + strconv.Itoa(cmds.PortWebs))

	http.HandleFunc("/", respHtml)
	http.HandleFunc("/text.json", putJson)
	http.HandleFunc("/favicon.ico", respStatic)
	http.HandleFunc("/css/", respStatic)
	http.HandleFunc("/purecss3/", respStatic)
	http.HandleFunc("/script/", respStatic)

	err := http.ListenAndServe(":"+strconv.Itoa(cmds.PortWebs), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
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

	fmt.Printf("RemoteAddr:\t%v\n", r.RemoteAddr)
	fmt.Printf("RequestURI:\t%v\n", r.RequestURI)

}
