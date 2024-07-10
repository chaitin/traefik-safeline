package detection

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/xbingW/traefik_safeline/sdk/pkg/misc"
)

type Request interface {
	Header() ([]byte, error)
	Body() (uint32, io.ReadCloser, error)
	Extra() ([]byte, error)
}

type HttpRequest struct {
	req *http.Request
	dc  *DetectionContext // this is optional
}

func MakeHttpRequest(req *http.Request) *HttpRequest {
	return &HttpRequest{
		req: req,
	}
}

func MakeHttpRequestInCtx(req *http.Request, dc *DetectionContext) *HttpRequest {
	ret := &HttpRequest{
		req: req,
		dc:  dc,
	}
	dc.Request = ret

	if dc.ReqBeginTime == 0 {
		dc.ReqBeginTime = misc.Now()
	}

	return ret
}

func (r *HttpRequest) Header() ([]byte, error) {
	var buf bytes.Buffer
	proto := r.req.Proto
	if r.dc != nil {
		if r.dc.Protocol != "" {
			proto = r.dc.Protocol
		} else {
			r.dc.Protocol = proto
		}
	}
	startLine := fmt.Sprintf("%s %s %s\r\n", r.req.Method, r.req.URL.RequestURI(), proto)
	_, err := buf.Write([]byte(startLine))
	if err != nil {
		return nil, err
	}
	_, err = buf.Write([]byte(fmt.Sprintf("Host: %s\r\n", r.req.Host)))
	if err != nil {
		return nil, err
	}
	err = r.req.Header.Write(&buf)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write([]byte("\r\n"))
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (r *HttpRequest) Body() (uint32, io.ReadCloser, error) {
	bodyBytes, err := ioutil.ReadAll(r.req.Body)
	if err != nil {
		return 0, nil, err
	}
	r.req.Body = ioutil.NopCloser(bytes.NewReader(bodyBytes))
	return uint32(len(bodyBytes)), ioutil.NopCloser(bytes.NewReader(bodyBytes)), nil
}

func (r *HttpRequest) Extra() ([]byte, error) {
	if r.dc == nil {
		return PlaceholderRequestExtra(misc.GenUUID()), nil
	}
	return GenRequestExtra(r.dc), nil
}
