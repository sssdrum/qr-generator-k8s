package main

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	"github.com/skip2/go-qrcode"
)

var db *sql.DB

type QRReq struct {
	URL string `json:"url"`
}

func initDB() {
	var err error

	// Get database connection info from environment variables
	dbHost := getEnv("DB_HOST", "postgres-service")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("POSTGRES_USER", "qruser")
	dbPassword := getEnv("POSTGRES_PASSWORD", "qrpassword123")
	dbName := getEnv("POSTGRES_DB", "qrdb")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	log.Println("Database connected successfully")

	// Create table if it doesn't exist
	createTable := `
	CREATE TABLE IF NOT EXISTS qr_cache (
		url_hash VARCHAR(32) PRIMARY KEY,
		url TEXT NOT NULL,
		qr_data BYTEA NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`

	_, err = db.Exec(createTable)
	if err != nil {
		log.Fatal("Failed to create table:", err)
	}

	log.Println("QR cache table ready")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func hashURL(url string) string {
	hash := md5.Sum([]byte(url))
	return hex.EncodeToString(hash[:])
}

func getCachedQR(urlHash string) ([]byte, bool) {
	var qrData []byte
	err := db.QueryRow("SELECT qr_data FROM qr_cache WHERE url_hash = $1", urlHash).Scan(&qrData)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, false
		}
		log.Printf("Error querying cache: %v", err)
		return nil, false
	}
	return qrData, true
}

func cacheQR(urlHash, url string, qrData []byte) {
	_, err := db.Exec("INSERT INTO qr_cache (url_hash, url, qr_data) VALUES ($1, $2, $3)",
		urlHash, url, qrData)
	if err != nil {
		log.Printf("Error caching QR: %v", err)
	}
}

func hello(w http.ResponseWriter, req *http.Request) {
	log.Printf("Hello endpoint accessed from %s", req.RemoteAddr)
	fmt.Fprint(w, "hello\n")
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
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	defer req.Body.Close()

	// Validate URL is provided
	if qrReq.URL == "" {
		log.Printf("Empty URL provided")
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	// Generate hash for the URL
	urlHash := hashURL(qrReq.URL)

	// Check cache first
	if cachedQR, found := getCachedQR(urlHash); found {
		log.Printf("QR code found in cache for URL: %s", qrReq.URL)
		w.Header().Set("Content-Type", "image/png")
		w.Header().Set("Content-Disposition", "inline; filename=qrcode.png")
		w.Header().Set("X-Cache", "HIT")
		w.Write(cachedQR)
		return
	}

	// Generate new QR code
	log.Printf("Generating new QR code for URL: %s", qrReq.URL)
	png, err := qrcode.Encode(qrReq.URL, qrcode.Medium, 256)
	if err != nil {
		log.Printf("QR generation failed: %v", err)
		http.Error(w, "Error generating QR code", http.StatusInternalServerError)
		return
	}

	// Cache the QR code
	cacheQR(urlHash, qrReq.URL, png)

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Disposition", "inline; filename=qrcode.png")
	w.Header().Set("X-Cache", "MISS")
	w.Write(png)

	log.Printf("QR code generated and cached for: %s", qrReq.URL)
}

func main() {
	// Initialize database
	initDB()
	defer db.Close()

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
