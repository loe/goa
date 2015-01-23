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

func main() {
	addr := flag.String("addr", ":8080", "proxy listen address")
	flag.Parse()

	director := func(req *http.Request) {
		auth := req.Header.Get("Authorization")
		log.Printf("Authorization: %s", auth)
	}

	proxy := AuthProxy{Director: director}
	log.Fatal(http.ListenAndServe(*addr, proxy))
}
