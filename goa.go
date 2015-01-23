package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

type AuthProxy struct {
	Director func(*http.Request)
}

func (p AuthProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.Director(r)
	fmt.Fprintln(w, "Great success!")
}

func logAuthorization(req *http.Request) {
	auth := req.Header.Get("Authorization")
	log.Printf("Authorization: %s", auth)
}

func main() {
	addr := flag.String("addr", ":8080", "proxy listen address")
	flag.Parse()

	proxy := AuthProxy{Director: logAuthorization}
	log.Fatal(http.ListenAndServe(*addr, proxy))
}
