package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

type Config struct {
	JssIP         string
	WebhookToPort map[string]int
}

var (
	Info  *log.Logger
	Error *log.Logger
)
var config Config
var jpwmux *http.ServeMux

func initLog(infoHandle io.Writer, errorHandle io.Writer) {
	Info = log.New(infoHandle,
		"INFO: ",
		log.Ldate|log.Ltime)
	Error = log.New(errorHandle,
		"ERROR: ",
		log.Ldate|log.Ltime)
}

// filterHandler sends a fordbidden response to requests that don't come from
// the ip specified in the config
func filterHandler(w http.ResponseWriter, r *http.Request) {
	if strings.Index(r.RemoteAddr, config.JssIP) != 0 {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(http.StatusText(http.StatusForbidden) + "\n"))
		return
	}

	jpwmux.ServeHTTP(w, r)
}

// registerHandlers takes the WebhookToPort map from the toml config and creates
// http handlers that forward requests from the JSS to the services on localhost
func registerHandlers(m *http.ServeMux, pmap map[string]int) {
	for i, j := range pmap {
		u, err := url.Parse(fmt.Sprintf("http://localhost:%d", j))
		if err != nil {
			Error.Fatal(err)
		}
		m.Handle(fmt.Sprintf("/%s", i), httputil.NewSingleHostReverseProxy(u))
		Info.Println("Forwarding", i, "to port", j)
	}
}

func main() {
	var err error
	initLog(os.Stdout, os.Stderr)
	_, err = toml.DecodeFile("config.toml", &config)
	if err != nil {
		Error.Fatal(err)
	}

	jpwmux = http.NewServeMux()
	registerHandlers(jpwmux, config.WebhookToPort)

	err = http.ListenAndServeTLS(":8443",
		"certs/server.crt",
		"certs/server.key",
		http.HandlerFunc(filterHandler))
	if err != nil {
		Error.Fatal(err)
	}
}
