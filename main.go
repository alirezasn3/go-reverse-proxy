package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func main() {
	backendRemote, err := url.Parse("http://localhost:8080")
	if err != nil {
		panic(err)
	}
	frontendRemote, err := url.Parse("http://localhost:3000")
	if err != nil {
		panic(err)
	}

	handler := func(p *httputil.ReverseProxy, remote *url.URL) func(http.ResponseWriter, *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			log.Println(r.URL)
			r.Host = remote.Host
			w.Header().Set("X-Ben", "Rad")
			p.ServeHTTP(w, r)
		}
	}

	backendProxy := httputil.NewSingleHostReverseProxy(backendRemote)
	frontendProxy := httputil.NewSingleHostReverseProxy(frontendRemote)
	http.HandleFunc("api.localhost/", handler(backendProxy, backendRemote))
	http.HandleFunc("localhost/", handler(frontendProxy, frontendRemote))

	err = http.ListenAndServe(":80", nil)
	if err != nil {
		panic(err)
	}
}
