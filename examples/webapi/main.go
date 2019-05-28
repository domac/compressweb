package main

import (
	. "../../../compressweb"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type Resp struct {
	Data string
	Seq  int
}

func CgiCmd(rspWriter http.ResponseWriter, req *http.Request) {

	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Println(err)
	}

	log.Printf("request data len = %d\n", len(data))

	//是否允许解压请求报文
	if ShouldUnCompress(req) {
		log.Println("should unCompress")
		data = GetUnCompressData(data)
	}

	rsp := &Resp{
		Data: string(data),
		Seq:  1024,
	}
	respData, _ := json.Marshal(rsp)

	//是否允许压缩响应报文
	if ShouldCompress(req) {
		log.Println("response using gzip")
		respData = GetCompressData(respData)
		SetHeader(rspWriter)
	}
	log.Printf("response len=%d\n", int(len(respData)))
	log.Printf("response data = %s \n",  respData)
	rspWriter.Header().Set("Content-Length", strconv.Itoa(int(len(respData))))
	_, _ = rspWriter.Write(respData)
}

func main() {
	http.HandleFunc("/cmd.cgi", CgiCmd)
	if err := http.ListenAndServe("0.0.0.0:28080", nil); err != nil {
		log.Fatal(err)
	}
}
