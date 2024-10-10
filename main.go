package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

func main() {
	var port int
	var listen string
	flag.IntVar(&port, "port", 4001, "Port to listen on")
	flag.StringVar(&listen, "listen", "localhost", "IP address to listen on")
	flag.Parse()

	args := flag.Args()
	if len(args) != 2 {
		fmt.Println("Usage: ./go-cors [-port <port>] [-listen <ip>] <allowed_origins> <backend_url>")
		os.Exit(1)
	}

	allowedOrigins := strings.Split(args[0], ",")
	backendURL := args[1]

	log.Printf("Allowed origins: %v", allowedOrigins)
	log.Printf("Backend URL: %s", backendURL)

	backend, err := url.Parse(backendURL)
	if err != nil {
		log.Fatal(err)
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
			log.Printf("Original response headers: %v", resp.Header)

			origin := resp.Request.Header.Get("Origin")
			resp.Header.Set("Access-Control-Allow-Origin", origin)
			resp.Header.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			resp.Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, x-tenant-id, x-request-id")
			resp.Header.Set("Access-Control-Allow-Credentials", "true")

			log.Printf("Modified response headers: %v", resp.Header)

			// Log response body for debugging
			bodyBytes, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			log.Printf("Response body: %s", bodyBytes)
			resp.Body = io.NopCloser(strings.NewReader(string(bodyBytes)))

			return nil
		},
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request: %s %s from %s", r.Method, r.URL, r.RemoteAddr)
		log.Printf("Request headers: %v", r.Header)

		origin := r.Header.Get("Origin")
		log.Printf("Origin header: %s", origin)

		if r.Method == "OPTIONS" {
			log.Println("Handling OPTIONS request")
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, x-tenant-id, x-request-id")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.WriteHeader(http.StatusOK)
			log.Printf("Responded to OPTIONS request. Headers: %v", w.Header())
			return
		}

		log.Println("Forwarding non-OPTIONS request to proxy handler")
		proxy.ServeHTTP(w, r)
	})

	listenAddr := fmt.Sprintf("%s:%d", listen, port)
	log.Printf("Starting proxy server on %s, forwarding to %s", listenAddr, backendURL)
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}