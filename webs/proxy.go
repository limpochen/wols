package webs

import (
	"bytes"
	"fmt"
	"net"
	"wols/config"
	"wols/llog"
)

// Start a proxy server listen on listenport
// This proxy will forward all HTTP request to httpport, and all HTTPS request to httpsport
func proxyStart() {
	proxylistener, err := net.Listen("tcp", fmt.Sprintf(":%d", config.Cfg.Webs.WebsPort))
	if err != nil {
		fmt.Printf("Unable to listen on: %d, error: %s/n", config.Cfg.Webs.WebsPort, err.Error())
	}
	defer proxylistener.Close()

	llog.Debug("proxy listen: " + fmt.Sprintf(":%d", config.Cfg.Webs.WebsPort))
	for {
		proxyconn, err := proxylistener.Accept()
		if err != nil {
			fmt.Printf("Unable to accept a request, error: %s/n", err.Error())
			continue
		}

		// Read a header firstly in case you could have opportunity to check request
		// whether to decline or proceed the request
		buffer := make([]byte, 1024)
		n, err := proxyconn.Read(buffer)
		if err != nil {
			//fmt.Printf("Unable to read from input, error: %s/n", err.Error())
			continue
		}

		var targetport int
		if isHTTPRequest(buffer) {
			targetport = config.HttpPort
			llog.Debug("proxy recive a http request")
		} else {
			targetport = config.HttpsPort
			llog.Debug("proxy recive a https request")
		}

		targetconn, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", targetport))
		if err != nil {
			fmt.Printf("Unable to connect to: %d, error: %s/n", targetport, err.Error())
			proxyconn.Close()
			continue
		}

		_, err = targetconn.Write(buffer[:n])
		if err != nil {
			fmt.Printf("Unable to write to output, error: %s/n", err.Error())
			proxyconn.Close()
			targetconn.Close()
			continue
		}

		go proxyRequest(proxyconn, targetconn)
		go proxyRequest(targetconn, proxyconn)
	}
}

// Forward all requests from r to w
func proxyRequest(r net.Conn, w net.Conn) {
	defer r.Close()
	defer w.Close()

	var buffer = make([]byte, 4096000)
	for {
		n, err := r.Read(buffer)
		if err != nil {
			//fmt.Printf("Unable to read from input, error: %s/n", err.Error())
			break
		}

		_, err = w.Write(buffer[:n])
		if err != nil {
			fmt.Printf("Unable to write to output, error: %s/n", err.Error())
			break
		}
	}
}

func isHTTPRequest(buffer []byte) bool {
	httpMethod := []string{"GET", "PUT", "HEAD", "POST", "DELETE", "PATCH", "OPTIONS"}
	for cnt := 0; cnt < len(httpMethod); cnt++ {
		if bytes.HasPrefix(buffer, []byte(httpMethod[cnt])) {
			return true
		}
	}
	return false
}
