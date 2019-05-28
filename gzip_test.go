package compressweb

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"
)

func TestGzipWriter(t *testing.T) {
	testData := strings.Repeat("apple", 100)
	ret := GetCompressData([]byte(testData))
	t.Logf("compress: size %d : %s!\n", len(ret), ret)
	ret = GetUnCompressData(ret)
	t.Logf("uncompress: size: %d, %s!\n", len(ret), ret)
}

func BenchmarkGzipWriter(b *testing.B) {
	testData := strings.Repeat("apple", 1024)
	for i := 0; i < b.N; i++ {
		GetCompressData([]byte(testData))
	}
}

func BenchmarkGzipReader(b *testing.B) {
	testData := strings.Repeat("apple", 1024)
	ret := GetCompressData([]byte(testData))

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		origData := GetUnCompressData(ret)
		if len(origData) != 5*1024 {
			b.Fail()
		}
	}
	b.StopTimer()
}

func BenchmarkNoUnGzip(b *testing.B) {
	testData := strings.Repeat("apple", 1024)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		bb := bytes.NewReader([]byte(testData))
		ret, _ := ioutil.ReadAll(bb)
		if len(ret) != 5*1024 {
			b.Fail()
		}
	}
	b.StopTimer()
}