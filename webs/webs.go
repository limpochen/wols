package webs

import (
	"fmt"
	"net"
	"net/http"
	"wols/config"
	"wols/llog"
)

func httpServ(ch chan string) {
	// pool := x509.NewCertPool()
	// caCrt, err := os.ReadFile(config.Cfg.Webs.RootCA)
	// if err != nil {
	// 	log.Fatalln("ReadFile err:", err)
	// }
	// pool.AppendCertsFromPEM([]byte(caCrt))

	var srv http.Server
	mux := http.NewServeMux()

	mux.HandleFunc("/opt/", basicAuth(respOpt))
	mux.HandleFunc("/recents", basicAuth(putJson))
	mux.HandleFunc("/", basicAuth(respStatic))

	if config.Cfg.Webs.EnableTls {
		srv.Addr = fmt.Sprintf("localhost:%d", config.HttpsPort)
	} else {
		srv.Addr = fmt.Sprintf(":%d", config.Cfg.Webs.WebsPort)
	}
	srv.Handler = mux
	// srv.TLSConfig = &tls.Config{
	// 	MinVersion: tls.VersionTLS12,
	// 	ClientCAs:  pool,
	// 	ClientAuth: tls.RequireAndVerifyClientCert,
	// }

	ln, err := net.Listen("tcp", srv.Addr)
	if err != nil {
		llog.Error("Web service: " + err.Error())
		ch <- "error"
		return
	}
	defer ln.Close()

	llog.Debug("http(s) listen: " + srv.Addr)
	ch <- "ok"

	if config.Cfg.Webs.EnableTls {
		err = srv.ServeTLS(ln, config.Cfg.Webs.CertFile, config.Cfg.Webs.KeyFile)
	} else {
		err = srv.Serve(ln)
	}
	if err != nil && err != http.ErrServerClosed {
		llog.Error("Web service: " + err.Error())
		ch <- "error"
	} else {
		ch <- "shutdown"
	}
}

func httpRedir() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// w.Header().Set("Cache-Control", "must-revalidate, no-store")
		// w.Header().Set("Content-Type", " text/html;charset=UTF-8")
		// w.Header().Set("Location", "https://"+r.Host)
		// w.WriteHeader(http.StatusTemporaryRedirect)
		http.Redirect(w, r, "https://"+r.Host, http.StatusTemporaryRedirect)
		llog.Debug("redir to: " + "https://" + r.Host)
	})

	llog.Debug("http redir listen: " + fmt.Sprintf("localhost:%d", config.HttpPort))
	http.ListenAndServe(fmt.Sprintf("localhost:%d", config.HttpPort), nil)
}

func WEBServ(ch chan string) {
	go httpServ(ch)
	if config.Cfg.Webs.EnableTls {
		go httpRedir()
		go proxyStart()
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
