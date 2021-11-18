package main

import (
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const defaultWaziappAddr = "unix:///var/lib/waziapp/proxy.sock"

var path string

func main() {
	log.SetFlags(0)

	var err error
	if len(os.Args) != 2 {
		usage()
		os.Exit(1)
	}
	switch os.Args[1] {
	case "-help", "--help", "-usage", "--usage":
		usage()
		os.Exit(0)
	}

	path = os.Args[1]
	if !strings.Contains(path, "://") {
		path = "http://" + path
	}
	if _, err := url.Parse(path); err != nil {
		log.Printf("Bad addr: %v", err)
		os.Exit(1)
	}

	addr := os.Getenv("WAZIAPP_ADDR")
	if addr == "" {
		addr = defaultWaziappAddr
	}
	i := strings.Index(addr, "://")
	if i == -1 {
		log.Printf("Bad WaziApp addr: Missing '://' in address: %s", addr)
		os.Exit(1)
	}

	listener, err := net.Listen(addr[:i], addr[i+3:])
	if err != nil {
		log.Printf("Listen error: %s", err.Error())
		os.Exit(1)
	}

	log.Printf("%s --> %s", addr, path)
	http.Serve(listener, http.HandlerFunc(handler))
}

func usage() {
	log.Printf("Usage: %s {addr}", os.Args[0])
	log.Printf("Create a unix domain socket (for use inside WaziApps) forwarding to local address.")
	log.Printf("Example:")
	log.Printf("  %s http://localhost:8080/test", os.Args[0])
	log.Printf("  unix:/var/lib/waziapp/proxy.sock --> http://localhost:8080/test")
	log.Printf("Use env WAZIAPP_ADDR to override the default waziapp proxy socket address.")
}

func handler(resp http.ResponseWriter, req *http.Request) {
	preoxyReq, err := http.NewRequest(req.Method, path+req.RequestURI, req.Body)
	if err != nil {
		log.Printf("%4s 502 %s %v", req.Method, req.RequestURI, err)
		resp.WriteHeader(http.StatusBadGateway)
		resp.Write([]byte(err.Error()))
		return
	}

	dup, err := http.DefaultClient.Do(preoxyReq)
	if err != nil {
		log.Printf("%4s 502 %s %v", req.Method, req.RequestURI, err)
		resp.WriteHeader(http.StatusBadGateway)
		resp.Write([]byte(err.Error()))
		return
	}
	for key, value := range dup.Header {
		resp.Header()[key] = value
	}
	log.Printf("%4s %d %s %v", req.Method, dup.StatusCode, req.RequestURI, err)
	resp.WriteHeader(dup.StatusCode)
	io.Copy(resp, dup.Body)
}
