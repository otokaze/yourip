package main

import (
	"encoding/json"
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
	resp   Response
	err    error
)

type Response struct {
	Code uint8  `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		ClientIP string `json:"client_ip"`
		ServerIP string `json:"server_ip"`
		Hostname string `json:"hostname"`
	} `json:"data"`
}

func init() {
	flag.BoolVar(&httpOn, "http", false, "Start HTTP Server.")
	flag.IntVar(&port, "p", 80, "The listening port.")
	flag.Parse()
}

func main() {
	if resp.Data.Hostname, err = os.Hostname(); err != nil {
		log.Printf("os.Hostname() error(%v)", err)
		return
	}
	var ipResp *http.Response
	if ipResp, err = http.Get("http://ifconfig.co/ip"); err != nil {
		log.Printf("http.Get(ifconfig.co/ip) error(%v)", err)
		return
	}
	var bs []byte
	if bs, err = ioutil.ReadAll(ipResp.Body); err != nil {
		log.Printf("ioutil.ReadAll() error(%v)", err)
		return
	}
	resp.Data.ServerIP = strings.TrimSpace(string(bs))
	data = []string{
		fmt.Sprintf("ServerIP: %s", resp.Data.ServerIP),
		fmt.Sprintf("Hostname: %s", resp.Data.Hostname),
	}
	if httpOn {
		http.HandleFunc("/", handleRoot)
		log.Fatalln(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
		return
	}
	fmt.Println(strings.Join(data, "\n"))
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var bs []byte
	if r.Form.Get("format") == "json" {
		resp.Data.ClientIP = r.RemoteAddr
		bs, _ = json.Marshal(&resp)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
	} else {
		clientIP := fmt.Sprintf("ClientIP: %s", r.RemoteAddr)
		bs = []byte(strings.Join(append([]string{clientIP}, data...), "\n"))
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	}
	if _, err = w.Write(bs); err != nil {
		log.Printf("resp.Write() error(%v)", err)
		return
	}
}
