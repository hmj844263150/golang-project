package main

import (
	"bytes"
	"compress/gzip"
	"espressif.com/chip/factory/api"
	"espressif.com/chip/factory/config"
	"espressif.com/chip/factory/db"
	"espressif.com/chip/factory/rpc"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

type heng struct{}

func (h *heng) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req, resp := rpc.NewRequest(), rpc.NewResponse()
	req.R, resp.W = r, w
	rpc.R2R(req, r)
	err := api.Router.Dispatch(req, resp)
	if err != nil {
		resp.Status = 500
		resp.Err = err
	}
	if resp.Err == rpc.Skip {
		return
	}
	output := resp.Json()
	w.Header().Set("Content-Type", "application/json")
	if len(output) > 32768 {
		acceptEncoding := r.Header.Get("Accept-Encoding")
		if strings.Contains(acceptEncoding, "gzip") {
			w.Header().Set("Content-Encoding", "gzip")
			output = gzipBytes(output)
		}
	}
	w.Write(output)
}

func gzipBytes(in []byte) []byte {
	bf := bytes.NewBuffer([]byte{})
	w := gzip.NewWriter(bf)
	w.Write(in)
	w.Close()
	out := bf.Bytes()
	return out
}

func main() {
	log.Println("heng...heng")
	defer db.Defer()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)

	if config.Cfg.Http != "" {
		s := &http.Server{
			Addr:           config.Cfg.Http,
			Handler:        &heng{},
			ReadTimeout:    60 * time.Second,
			WriteTimeout:   60 * time.Second,
			MaxHeaderBytes: 10 * 1024 * 1024,
		}
		go func() {
			log.Fatal(s.ListenAndServe())
		}()
	}
	if config.Cfg.TlsHttp != "" && config.Cfg.TlsConfig != nil {
		s := &http.Server{
			Addr:           config.Cfg.TlsHttp,
			Handler:        &heng{},
			ReadTimeout:    60 * time.Second,
			WriteTimeout:   60 * time.Second,
			MaxHeaderBytes: 10 * 1024 * 1024,
		}
		go func() {
			log.Fatal(s.ListenAndServeTLS(config.Cfg.CertFile, config.Cfg.KeyFile))
		}()
	}

	<-sc
	log.Println("quit heng...")
}
