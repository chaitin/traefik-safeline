package t1k

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/xbingW/traefik_safeline/sdk/pkg/misc"
)

func TestWriteSimpleSection(t *testing.T) {
	sec := MakeSimpleSection(TAG_VERSION, []byte("Proto:2\n"))
	var buf bytes.Buffer
	err := WriteSection(sec, &buf)
	if err != nil {
		t.Fatal(err)
	}
	misc.PrintHex(buf.Bytes())
}

func TestWriteReaderSection(t *testing.T) {
	rbuf := bytes.NewBuffer([]byte(
		"GET /webshell.php HTTP/1.1\r\n" +
			"Host: a.com\r\n" +
			"Content-Length: 65533\r\n\r\n",
	))
	sec := MakeReaderSection(TAG_HEADER, uint32(len(rbuf.Bytes())), rbuf)
	var buf bytes.Buffer
	err := WriteSection(sec, &buf)
	if err != nil {
		t.Fatal(err)
	}
	misc.PrintHex(buf.Bytes())
}

func TestSectionSplit(t *testing.T) {
	data := []byte{
		0x01, 0x03, 0x00, 0x00, 0x00, // first section -- tag=1, len=3
		0xaa, 0xaa, 0xaa,
		0x02, 0x05, 0x00, 0x00, 0x00, // second section -- tag=2, len=5
		0xbb, 0xbb, 0xbb, 0xbb, 0xbb,
		0x03, 0x07, 0x00, 0x00, 0x00, // third section -- tag=3, len=7
		0xcc, 0xcc, 0xcc, 0xcc, 0xcc, 0xcc, 0xcc,
	}
	buf := bytes.NewBuffer(data)
	sec1, err := ReadFullSection(buf)
	if err != nil {
		t.Fatal(err)
	}
	sec2, err := ReadSection(buf)
	if err != nil {
		t.Fatal(err)
	}
	var bodyBuf bytes.Buffer
	err = sec2.WriteBody(&bodyBuf)
	if err != nil {
		t.Fatal(err)
	}
	sec3, err := ReadFullSection(buf)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%v\n%v %v\n%v\n", sec1, sec2, bodyBuf, sec3)
}
