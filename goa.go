package main

import (
	"flag"
	"github.com/elazarl/goproxy"
	"log"
	"net/http"
)

func main() {
	verbose := flag.Bool("v", false, "should every proxy request be logged to stdout")
	addr := flag.String("addr", ":8080", "proxy listen address")
	flag.Parse()
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = *verbose

	proxy.OnRequest().DoFunc(
		func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			auth := r.Header.Get("Authorization")
			log.Printf("Authorization: %s", auth)
			return r, nil
		})

	log.Fatal(http.ListenAndServe(*addr, proxy))
}
