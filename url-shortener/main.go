package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"sync"
)

var (
	urlStore = make(map[string]string)
	mutex    sync.Mutex
)

func main() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/shorten", shortenURL)
	http.HandleFunc("/r/", redirect)
	http.HandleFunc("/health", healthCheck)

	fmt.Println("URL Shortener running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func homePage(w http.ResponseWriter, r *http.Request) {
	html := `
	<h1>URL Shortener</h1>
	<form action="/shorten" method="POST">
		<input type="text" name="url" placeholder="Enter URL here" size="40"/>
		<button type="submit">Shorten</button>
	</form>
	`
	fmt.Fprintf(w, html)
}

func shortenURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	originalURL := r.FormValue("url")
	if originalURL == "" {
		http.Error(w, "URL cannot be empty", http.StatusBadRequest)
		return
	}

	if !strings.HasPrefix(originalURL, "http") {
		originalURL = "https://" + originalURL
	}

	code := generateCode(6)

	mutex.Lock()
	urlStore[code] = originalURL
	mutex.Unlock()

	shortURL := fmt.Sprintf("http://localhost:8080/r/%s", code)
	fmt.Fprintf(w, `
	<h1>URL Shortened!</h1>
	<p>Original: %s</p>
	<p>Shortened: <a href="%s">%s</a></p>
	<br/>
	<a href="/">Shorten another</a>
	`, originalURL, shortURL, shortURL)
}

func redirect(w http.ResponseWriter, r *http.Request) {
	code := strings.TrimPrefix(r.URL.Path, "/r/")

	mutex.Lock()
	originalURL, exists := urlStore[code]
	mutex.Unlock()

	if !exists {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusMovedPermanently)
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "OK")
}

func generateCode(length int) string {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	code := make([]byte, length)
	for i := range code {
		code[i] = chars[rand.Intn(len(chars))]
	}
	return string(code)
}