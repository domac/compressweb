package compressweb

import (
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
)

const (
	vary            = "Vary"
	AcceptEncoding  = "Accept-Encoding"
	ContentEncoding = "Content-Encoding"
	ContentLength   = "Content-Length"
)

var zippers = sync.Pool{}
var unZippers = sync.Pool{}

//初始化
func init() {

	log.Println("compress service init")

	zippers = sync.Pool{
		New: func() interface{} {
			wr := gzip.NewWriter(nil)
			return &CompressWriter{gzipWriter{wr}}
		}}

	unZippers = sync.Pool{New: func() interface{} {
		buf := []byte{31, 139, 8, 0, 0, 0, 0, 0, 0, 255}
		rbuf := bytes.NewReader(buf)
		wr, err := gzip.NewReader(rbuf)
		return &CompressReader{gzipReader{wr}, err}
	}}
}

///////////// CompressWriter /////////////////
type gzipWriter struct {
	*gzip.Writer
}

type CompressWriter struct {
	gzipWriter
}

func NewCompressWriter(w io.Writer) *CompressWriter {
	wr := zippers.Get().(*CompressWriter)
	wr.Reset(w)
	return wr
}

func (gw *CompressWriter) Close() error {
	e := gw.gzipWriter.Close()
	zippers.Put(gw)
	return e
}

///////////// CompressReader /////////////////
type gzipReader struct {
	*gzip.Reader
}

type CompressReader struct {
	gzipReader
	e error
}

func NewCompressReader(r io.Reader) (*CompressReader, error) {
	gr := unZippers.Get().(*CompressReader)
	if gr.e != nil {
		return nil, gr.e
	}
	_ = gr.Reset(r)
	return gr, nil
}

func (gr *CompressReader) Close() error {
	e := gr.gzipReader.Close()
	unZippers.Put(gr)
	return e
}

// Fix: https://github.com/mholt/caddy/issues/38
//shouldCompress时，要更新相关header
func SetHeader(rspWriter http.ResponseWriter) {
	rspWriter.Header().Set(ContentEncoding, "gzip")
	rspWriter.Header().Add(vary, AcceptEncoding)
	rspWriter.Header().Del(ContentLength)
	rspWriter.Header().Set("compressBy", "conn_sc")
}

func ShouldCompress(req *http.Request) bool {
	acceptGzip := false
	hdr := req.Header
	for _, encoding := range strings.Split(hdr.Get(AcceptEncoding), ",") {
		if "gzip" == strings.TrimSpace(encoding) {
			acceptGzip = true
			break
		}
	}
	return acceptGzip
}

func ShouldUnCompress(req *http.Request) bool {
	acceptGzip := false
	for _, encoding := range strings.Split(req.Header.Get(ContentEncoding), ",") {
		if "gzip" == strings.TrimSpace(encoding) {
			acceptGzip = true
			break
		}
	}
	return acceptGzip
}

func GetUnCompressData(data []byte) []byte {
	bc := bytes.NewReader(data)
	gzipReader, _ := NewCompressReader(bc)
	ret, _ := ioutil.ReadAll(gzipReader)
	_ = gzipReader.Close()
	return ret
}

func GetCompressData(data []byte) []byte {
	var ret bytes.Buffer
	zipper := NewCompressWriter(&ret)
	_, _ = zipper.Write(data)
	_ = zipper.Close()
	return ret.Bytes()
}
