package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
)

var config Config

type Proxy struct {
	Listen  string `json:"listen"`
	Connect string `json:"connect"`
}

type Config struct {
	Listen  string  `json:"listen"`
	HTTPS   bool    `json:"https"`
	Cert    string  `json:"cert"`
	Key     string  `json:"key"`
	Proxies []Proxy `json:"proxies"`
}

func init() {
	// Read config file
	execPath, err := os.Executable()
	if err != nil {
		panic(err)
	}
	path := filepath.Dir(execPath)
	bytes, err := os.ReadFile(filepath.Join(path, "config.json"))
	if err != nil {
		panic(err)
	}

	// Parse config file into global config variable
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		panic(err)
	}
}

func main() {
	// Loop over proxies array
	for _, proxy := range config.Proxies {
		go func(c string, l string) {
			// Handle requests on the url that the proxy should listen on
			http.HandleFunc(l, func(w http.ResponseWriter, r *http.Request) {
				// Parse the url that the proxy should forward request to
				connectUrl, err := url.Parse(c)
				if err != nil {
					panic(err)
				}
				// Create reverse proxy from url
				reverseProxy := httputil.NewSingleHostReverseProxy(connectUrl)
				// Log request info
				log.Printf("[%s] -> [%s]\n", r.URL, c)
				// Change request's host to destination host
				r.Host = connectUrl.Host
				// Forward request to destination
				reverseProxy.ServeHTTP(w, r)
			})
		}(proxy.Connect, proxy.Listen)
	}

	// Create http server and listen for requests
	if config.HTTPS {
		if err := http.ListenAndServeTLS(config.Listen, config.Cert, config.Key, nil); err != nil {
			panic(err)
		}
	} else {
		if err := http.ListenAndServe(config.Listen, nil); err != nil {
			panic(err)
		}
	}
}
