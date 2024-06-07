package traefik_safeline

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/xbingW/t1k/detector"
)

// Package example a example plugin.

// Config the plugin configuration.
type Config struct {
	// Addr is the address for the detector
	Addr string `json:"addr"`
	// Get ip from header, if not set, get ip from remote addr
	IpHeader string `json:"ip_header"`
	// When ip_header has multiple ip, use this to get the last ip
	//
	//for example, X-Forwarded-For: ip1, ip2, ip3
	// 	when ip_last_index is 0, the client ip is ip3
	// 	when ip_last_index is 1, the client ip is ip2
	// 	when ip_last_index is 2, the client ip is ip1
	IpLastIndex uint `json:"ip_last_index"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		Addr:        "",
		IpHeader:    "",
		IpLastIndex: 0,
	}
}

// Safeline a plugin.
type Safeline struct {
	next   http.Handler
	name   string
	config *Config
	logger *log.Logger
}

// New created a new plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	return &Safeline{
		next:   next,
		name:   name,
		config: config,
		logger: log.New(os.Stdout, "safeline", log.LstdFlags),
	}, nil
}

func (s *Safeline) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	d := detector.NewDetector(detector.Config{
		Addr:        s.config.Addr,
		IpHeader:    s.config.IpHeader,
		IpLastIndex: s.config.IpLastIndex,
	})
	resp, err := d.DetectorRequest(req)
	if err != nil {
		s.logger.Printf("Failed to detect request: %v", err)
	}
	if resp != nil && !resp.Allowed() {
		rw.WriteHeader(resp.StatusCode())
		if err := json.NewEncoder(rw).Encode(resp.BlockMessage()); err != nil {
			s.logger.Printf("Failed to encode block message: %v", err)
		}
		return
	}
	s.next.ServeHTTP(rw, req)
}
