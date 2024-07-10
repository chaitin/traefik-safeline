package traefik_safeline

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/xbingW/t1k"
)

// Package example a example plugin.

// Config the plugin configuration.
type Config struct {
	// Addr is the address for the detector
	Addr     string `yaml:"addr"`
	PoolSize int    `yaml:"pool_size"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		Addr:     "",
		PoolSize: 100,
	}
}

// Safeline a plugin.
type Safeline struct {
	next   http.Handler
	server *t1k.Server
	name   string
	config *Config
	logger *log.Logger
}

// New created a new plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	logger := log.New(os.Stdout, "safeline", log.LstdFlags)
	logger.Printf("config: %+v", config)
	server, err := t1k.NewWithPoolSize(config.Addr, config.PoolSize)
	if err != nil {
		return nil, err
	}
	return &Safeline{
		next:   next,
		name:   name,
		config: config,
		server: server,
		logger: logger,
	}, nil
}

func (s *Safeline) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			s.logger.Printf("panic: %s", r)
		}
	}()
	rw.Header().Set("X-Chaitin-waf", "safeline")
	result, err := s.server.DetectHttpRequest(req)
	if err != nil {
		s.logger.Printf("error in detection: \n%+v\n", err)
		s.next.ServeHTTP(rw, req)
		return
	}
	if result.Blocked() {
		rw.WriteHeader(http.StatusForbidden)
		_, _ = rw.Write([]byte("Blocked by safeline\n"))
		return
	}
	s.next.ServeHTTP(rw, req)
	//rw.WriteHeader(http.StatusForbidden)
	//_, _ = rw.Write([]byte("Inject by safeline\n"))
}
