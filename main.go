package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

var (
	httpOn bool
	port   int
	data   []string
	err    error
)

func init() {
	flag.BoolVar(&httpOn, "http", false, "Start HTTP Server.")
	flag.IntVar(&port, "p", 80, "The listening port.")
	flag.Parse()
}

func main() {
	var host string
	if host, err = os.Hostname(); err != nil {
		log.Printf("os.Hostname() error(%v)", err)
		return
	}
	var resp *http.Response
	if resp, err = http.Get("http://ifconfig.co/ip"); err != nil {
		log.Printf("http.Get(ifconfig.co/ip) error(%v)", err)
		return
	}
	var bs []byte
	if bs, err = ioutil.ReadAll(resp.Body); err != nil {
		log.Printf("ioutil.ReadAll() error(%v)", err)
		return
	}
	data = []string{
		fmt.Sprintf("ServerIP: %s", strings.TrimSpace(string(bs))),
		fmt.Sprintf("Hostname: %s", host),
	}
	if httpOn {
		http.HandleFunc("/", handleRoot)
		log.Fatalln(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
		return
	}
	fmt.Println(strings.Join(data, "\n"))
}

func handleRoot(resp http.ResponseWriter, req *http.Request) {
	var data = append([]string{fmt.Sprintf("ClientIP: %s", req.RemoteAddr)}, data...)
	resp.Header().Set("Content-Type", "text/plain; charset=utf-8")
	if _, err = resp.Write([]byte(strings.Join(data, "\n"))); err != nil {
		log.Printf("resp.Write() error(%v)", err)
		return
	}
}
