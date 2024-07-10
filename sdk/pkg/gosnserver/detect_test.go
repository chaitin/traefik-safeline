package gosnserver

import (
	"bufio"
	"bytes"
	"fmt"
	"net/http"
	"testing"

	"github.com/xbingW/traefik_safeline/sdk/pkg/detection"
	"github.com/xbingW/traefik_safeline/sdk/pkg/misc"
)

func TestWriteDetectRequest(t *testing.T) {
	sReq := "POST /form.php?id=3 HTTP/1.1\r\n" +
		"Host: a.com\r\n" +
		"Content-Length: 40\r\n" +
		"Content-Type: application/json\r\n\r\n" +
		"{\"name\": \"youcai\", \"password\": \"******\"}"
	req, err := http.ReadRequest(bufio.NewReader(bytes.NewBuffer([]byte(sReq))))
	if err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	err = writeDetectionRequest(&buf, detection.MakeHttpRequest(req))
	if err != nil {
		t.Fatal(err)
	}
	misc.PrintHex(buf.Bytes())
}

func TestWriteDetectRequestAndResponse(t *testing.T) {
	sReq := "POST /form.php HTTP/1.1\r\n" +
		"Host: a.com\r\n" +
		"Content-Length: 40\r\n" +
		"Content-Type: application/json\r\n\r\n" +
		"{\"name\": \"youcai\", \"password\": \"******\"}"
	req, err := http.ReadRequest(bufio.NewReader(bytes.NewBuffer([]byte(sReq))))
	if err != nil {
		t.Fatal(err)
	}

	sRsp := "HTTP/1.1 200 OK\r\n" +
		"Content-Length: 29\r\n" +
		"Content-Type: application/json\r\n\r\n" +
		"{\"err\": \"password-incorrect\"}"
	rsp, err := http.ReadResponse(bufio.NewReader(bytes.NewBuffer([]byte(sRsp))), req)
	if err != nil {
		t.Fatal(err)
	}

	dc := detection.New()
	dreq := detection.MakeHttpRequestInCtx(req, dc)
	drsp := detection.MakeHttpResponseInCtx(rsp, dc)

	var buf bytes.Buffer
	err = writeDetectionRequest(&buf, dreq)
	if err != nil {
		t.Fatal(err)
	}
	dc.T1KContext = []byte("sample-t1k-context")
	err = writeDetectionResponse(&buf, drsp)
	if err != nil {
		t.Fatal(err)
	}

	misc.PrintHex(buf.Bytes())
}

func TestReadDetectResult(t *testing.T) {
	data := []byte{
		0x41, 0x01, 0x00, 0x00, 0x00, 0x2e, 0xa5, 0x4d,
		0x00, 0x00, 0x00, 0x7b, 0x22, 0x65, 0x76, 0x65,
		0x6e, 0x74, 0x5f, 0x69, 0x64, 0x22, 0x3a, 0x22,
		0x38, 0x36, 0x63, 0x38, 0x33, 0x62, 0x35, 0x61,
		0x33, 0x66, 0x62, 0x32, 0x34, 0x31, 0x61, 0x32,
		0x38, 0x39, 0x39, 0x37, 0x64, 0x39, 0x34, 0x65,
		0x34, 0x62, 0x32, 0x39, 0x63, 0x61, 0x65, 0x33,
		0x22, 0x2c, 0x22, 0x72, 0x65, 0x71, 0x75, 0x65,
		0x73, 0x74, 0x5f, 0x68, 0x69, 0x74, 0x5f, 0x77,
		0x68, 0x69, 0x74, 0x65, 0x6c, 0x69, 0x73, 0x74,
		0x22, 0x3a, 0x66, 0x61, 0x6c, 0x73, 0x65, 0x7d,
	}
	rBuf := bytes.NewBuffer(data)
	ret, err := readDetectionResult(rBuf)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%v\n", ret)
}
