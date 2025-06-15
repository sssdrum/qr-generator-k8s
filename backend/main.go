package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/skip2/go-qrcode"
)

func hello(w http.ResponseWriter, req *http.Request) {
	log.Printf("Hello endpoint accessed from %s", req.RemoteAddr)
	fmt.Fprint(w, "hello\n")
}

type QRReq struct {
	URL string `json:"url"`
}

func genQR(w http.ResponseWriter, req *http.Request) {
	log.Printf("QR generation request from %s", req.RemoteAddr)

	if req.Method != http.MethodPost {
		log.Printf("Invalid method %s attempted", req.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var qrReq QRReq

	err := json.NewDecoder(req.Body).Decode(&qrReq)
	if err != nil {
		log.Printf("JSON decode error: %v", err)
		fmt.Fprint(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	defer req.Body.Close()

	// Validate URL is provided
	if qrReq.URL == "" {
		log.Printf("Empty URL provided")
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	log.Printf("Generating QR code for URL: %s", qrReq.URL)

	var png []byte
	png, err = qrcode.Encode(qrReq.URL, qrcode.Medium, 256)
	if err != nil {
		log.Printf("QR generation failed: %v", err)
		http.Error(w, "Error generating qr code for url: %s", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Disposition", "inline; filename=qrcode.png")
	w.Write(png)

	log.Printf("QR code generated successfully for: %s", qrReq.URL)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/hello", hello)
	mux.HandleFunc("/api/generate-qr", genQR)

	loggedMux := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// Enable CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight requests
		if req.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		log.Printf("%s %s from %s", req.Method, req.URL.Path, req.RemoteAddr)
		mux.ServeHTTP(w, req)
	})

	log.Printf("Starting server...")
	log.Printf("Server listening at http://localhost:8080/")
	http.ListenAndServe(":8080", loggedMux)
}
