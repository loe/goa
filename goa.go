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
	auth := r.Header.Get("Authorization")
	log.Printf("Authorization: %s", auth)

	fmt.Fprintln(w, "Great success!")
}

func main() {
	addr := flag.String("addr", ":8080", "proxy listen address")
	flag.Parse()

	var proxy AuthProxy
	log.Fatal(http.ListenAndServe(*addr, proxy))
}
