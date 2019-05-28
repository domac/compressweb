package main

import (
	. "../../../compressweb"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var a = flag.Bool("a", false, "Accept-Gzip")
var c = flag.Bool("c", false, "Content-Gzip")
var p = flag.String("p", "", "post data")

func main() {
	flag.Parse()
	addr := "http://localhost:28080/cmd.cgi"
	fmt.Printf("addr = %s\n\n", addr)

	reqData := []byte(*p)
	client := &http.Client{
		Transport: &http.Transport{
			DisableCompression: true,
		},
	}
	var req *http.Request
	if *c {
		reqData = GetCompressData(reqData)
		req, _ = http.NewRequest("POST", addr, bytes.NewBuffer(reqData))
		req.Header.Set(ContentEncoding, "gzip")
		fmt.Printf("compress data to remote : len=%d\n", len(reqData))
	} else {
		req, _ = http.NewRequest("POST", addr, bytes.NewBuffer(reqData))
		fmt.Printf("data to remote : len=%d\n", len(reqData))
	}

	if *a {
		req.Header.Set(AcceptEncoding, "gzip")
	}
	//http 请求
	doReq(req, client)
}

func doReq(req *http.Request, client *http.Client) {
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer func() {
		if resp != nil && resp.Body != nil {
			_ = resp.Body.Close()
		}
	}()
	respData, err := ioutil.ReadAll(resp.Body)
	log.Printf("date origin length=%d\n", len(respData))
	ce := resp.Header.Get(ContentEncoding)
	if strings.Contains(ce, "gzip") {
		log.Println("response using gzip")
		respData = GetUnCompressData(respData)
	} else {
		log.Println("response using default")
	}
	fmt.Printf("data from remote : len=%d,  body=%s\n", len(respData), respData)
}
