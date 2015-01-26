package main

import (
	"flag"
	"fmt"
	"github.com/bradfitz/gomemcache/memcache"
	"log"
	"net/http"
)

type AuthProxy struct {
	Memcache *memcache.Client
}

func (proxy AuthProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	count, err := countAuthorization(r, proxy.Memcache)

	if err == nil {
		fmt.Fprintf(w, "Great success! => %d\n", count)
	} else {
		fmt.Fprintln(w, "Fail.")
	}
}

func countAuthorization(req *http.Request, m *memcache.Client) (count uint64, err error) {
	auth, _ := extractAuthorization(req)

	// attempt to increment by 1
	count, err = m.Increment(auth, 1)

	// if we get a miss perform an add, then incr again. Its fine if the add fails, it means another writer was first, the incr is still atomic.
	if err == memcache.ErrCacheMiss {
		// intentionally ignoring the error
		m.Add(&memcache.Item{Key: auth, Value: []byte("0")})
		count, err = m.Increment(auth, 1)
	}

	return count, err
}

func extractAuthorization(req *http.Request) (auth string, err error) {
	auth = req.Header.Get("Authorization")

	return auth, nil
}

func main() {
	addr := flag.String("addr", ":8080", "proxy listen address")
	flag.Parse()

	mc := memcache.New("localhost:11211")
	proxy := AuthProxy{Memcache: mc}
	log.Fatal(http.ListenAndServe(*addr, proxy))
}
