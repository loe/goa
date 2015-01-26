package main

import (
	"flag"
	"fmt"
	"github.com/bradfitz/gomemcache/memcache"
	"log"
	"net/http"
)

type tokenProxy struct {
	Memcache *memcache.Client
}

func (proxy tokenProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	count, err := countRequest(r, proxy.Memcache)

	if err == nil {
		fmt.Fprintf(w, "Great success! => %d\n", count)
	} else {
		fmt.Fprintln(w, "Fail.")
	}
}

func countRequest(req *http.Request, m *memcache.Client) (count uint64, err error) {
	token, _ := extractAuthorization(req)

	// attempt to increment by 1
	count, err = m.Increment(token, 1)

	// if we get a miss perform an add, then incr again. Its fine if the add fails, it means another writer was first, the incr is still atomic.
	if err == memcache.ErrCacheMiss {
		// intentionally ignoring the error
		m.Add(&memcache.Item{Key: token, Value: []byte("0")})
		count, err = m.Increment(token, 1)
	}

	return count, err
}

func extractAuthorization(req *http.Request) (token string, err error) {
	token = req.Header.Get("Authorization")

	return token, nil
}

func main() {
	addr := flag.String("addr", ":8080", "proxy listen address")
	flag.Parse()

	mc := memcache.New("localhost:11211")
	proxy := tokenProxy{Memcache: mc}
	log.Fatal(http.ListenAndServe(*addr, proxy))
}
