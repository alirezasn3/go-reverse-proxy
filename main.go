package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

var config Config

type Proxy struct {
	Listen  string `json:"listen"`
	Connect string `json:"connect"`
}

type Config struct {
	Listen  string  `json:"listen"`
	Proxies []Proxy `json:"proxies"`
}

// Load config file
func init() {
	// Default config file path is ./config.json
	configPath := "config.json"

	// Use the first command line argument as config file path if one is provided
	if len(os.Args) > 1 {
		configPath = os.Args[1] + configPath
	}

	// Read config file
	bytes, err := os.ReadFile(configPath)
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
	if err := http.ListenAndServe(config.Listen, nil); err != nil {
		panic(err)
	}
}
