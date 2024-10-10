package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"
)

type Service struct {
	Path    string `yaml:"path"`
	Backend string `yaml:"backend"`
}

type Config struct {
	Services []Service `yaml:"services"`
}

func main() {
	var port int
	var listen, configFile, corsOrigins string
	flag.IntVar(&port, "port", 8080, "Port to listen on")
	flag.StringVar(&listen, "listen", "127.0.0.1", "IP address to listen on")
	flag.StringVar(&configFile, "config", "", "Path to YAML config file")
	flag.StringVar(&corsOrigins, "cors", "", "Allowed CORS origins (comma-separated)")
	flag.Parse()

	if corsOrigins == "" {
		fmt.Println("Error: CORS origins must be specified")
		flag.Usage()
		os.Exit(1)
	}

	allowedOrigins := strings.Split(corsOrigins, ",")
	log.Printf("Allowed origins: %v", allowedOrigins)

	var config Config
	if configFile != "" {
		data, err := os.ReadFile(configFile)
		if err != nil {
			log.Fatalf("Error reading config file: %v", err)
		}
		err = yaml.Unmarshal(data, &config)
		if err != nil {
			log.Fatalf("Error parsing config file: %v", err)
		}
	} else if backend := flag.Arg(0); backend != "" {
		config.Services = []Service{{Path: "/", Backend: backend}}
	} else {
		fmt.Println("Error: Either -config or a backend URL must be provided")
		flag.Usage()
		os.Exit(1)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request: %s %s from %s", r.Method, r.URL, r.RemoteAddr)
		log.Printf("Request headers: %v", r.Header)

		origin := r.Header.Get("Origin")
		log.Printf("Origin header: %s", origin)

		if r.Method == "OPTIONS" {
			handleCORS(w, r, allowedOrigins)
			return
		}

		// First, check for exact path matches
		for _, service := range config.Services {
			if service.Path == r.URL.Path {
				handleProxy(w, r, service.Backend, allowedOrigins)
				return
			}
		}

		// If no exact match, check for regex matches
		for _, service := range config.Services {
			if match, _ := regexp.MatchString("^"+service.Path, r.URL.Path); match {
				handleProxy(w, r, service.Backend, allowedOrigins)
				return
			}
		}

		http.NotFound(w, r)
	})

	listenAddr := fmt.Sprintf("%s:%d", listen, port)
	log.Printf("Starting proxy server on %s", listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}

func handleCORS(w http.ResponseWriter, r *http.Request, allowedOrigins []string) {
	origin := r.Header.Get("Origin")
	if isAllowedOrigin(origin, allowedOrigins) {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, x-tenant-id, x-request-id")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
	}
	w.WriteHeader(http.StatusOK)
}

func handleProxy(w http.ResponseWriter, r *http.Request, backendURL string, allowedOrigins []string) {
	backend, err := url.Parse(backendURL)
	if err != nil {
		log.Printf("Error parsing backend URL: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	proxy := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.URL.Scheme = backend.Scheme
			req.URL.Host = backend.Host
			req.Host = backend.Host
			log.Printf("Forwarding request to backend: %s %s", req.Method, req.URL)
		},
		ModifyResponse: func(resp *http.Response) error {
			log.Printf("Received response from backend. Status: %s", resp.Status)
			origin := resp.Request.Header.Get("Origin")
			if isAllowedOrigin(origin, allowedOrigins) {
				resp.Header.Set("Access-Control-Allow-Origin", origin)
				resp.Header.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
				resp.Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, x-tenant-id, x-request-id")
				resp.Header.Set("Access-Control-Allow-Credentials", "true")
			}
			return nil
		},
	}

	proxy.ServeHTTP(w, r)
}

func isAllowedOrigin(origin string, allowedOrigins []string) bool {
	for _, allowed := range allowedOrigins {
		if allowed == "*" || allowed == origin {
			return true
		}
	}
	return false
}