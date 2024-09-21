package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"slices"

	goSystemd "github.com/alirezasn3/go-systemd"
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
	execPath, err := os.Executable()
	if err != nil {
		panic(err)
	}

	// check for install and uninstall commands
	if slices.Contains(os.Args, "--install") {
		err = goSystemd.CreateService(&goSystemd.Service{Name: "go-reverse-proxy", ExecStart: execPath, Restart: "on-failure", RestartSec: "3s"})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		} else {
			fmt.Println("go-reverse-proxy service created")
			os.Exit(0)
		}
	} else if slices.Contains(os.Args, "--uninstall") {
		err := goSystemd.DeleteService("go-reverse-proxy")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		} else {
			fmt.Println("go-reverse-proxy service deleted")
			os.Exit(0)
		}
	}

	// Read config file
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
