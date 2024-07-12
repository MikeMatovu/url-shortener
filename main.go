package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type UrlShortener struct {
	urls map[string]string
}

func (us *UrlShortener) HandleRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	// Render HTML or redirect to a specific page
	w.Header().Set("Content-Type", "text/html")
	responseHTML := `
        <h1>Welcome to the URL Shortener</h1>
        <form method="post" action="/shorten">
            <input type="text" name="url" placeholder="Enter a URL to shorten">
            <input type="submit" value="Shorten">
        </form>
    `
	fmt.Fprint(w, responseHTML)
}

func (us *UrlShortener) HandleShorten(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	originalURL := r.FormValue("url")

	if originalURL == "" {
		http.Error(w, "Url is required", http.StatusBadRequest)
		return
	}

	shortKey := generateShortKey()

	us.urls[shortKey] = originalURL

	shortenedURL := fmt.Sprintf("http://localhost:8080/short/%s", shortKey)

	//Render html with the shortened url
	w.Header().Set("Content-Type", "text/html")
	responseHTML := fmt.Sprintf(`
        <h2>URL Shortener</h2>
        <p>Original URL: %s</p>
        <p>Shortened URL: <a href="%s">%s</a></p>
        <form method="post" action="/shorten">
            <input type="text" name="url" placeholder="Enter a URL">
            <input type="submit" value="Shorten">
        </form>
    `, originalURL, shortenedURL, shortenedURL)
	fmt.Fprint(w, responseHTML)
}

func (us *UrlShortener) HandleRedirect(w http.ResponseWriter, r *http.Request) {
	shortKey := r.URL.Path[len("/short/"):]

	if shortKey == "" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	originalURL, ok := us.urls[shortKey]

	if !ok {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusFound)
}

func generateShortKey() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const keyLength = 6

	rand.Seed(time.Now().UnixNano())
	shortKey := make([]byte, keyLength)

	for i := range shortKey {
		shortKey[i] = charset[rand.Intn(len(charset))]
	}
	return string(shortKey)
}

func (us *UrlShortener) SaveToFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := gob.NewEncoder(file)
	return encoder.Encode(us.urls)
}

func (us *UrlShortener) LoadFromFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	decoder := gob.NewDecoder(file)
	return decoder.Decode(&us.urls)
}

func main() {
	urlShortener := UrlShortener{
		urls: make(map[string]string),
	}

	err := urlShortener.LoadFromFile("urls.gob")
	if err != nil {
		fmt.Println("No existing URLs found or error loading URLs:", err)
	}

	http.HandleFunc("/", urlShortener.HandleRoot)
	http.HandleFunc("/shorten", urlShortener.HandleShorten)
	http.HandleFunc("/short/", urlShortener.HandleRedirect)

	fmt.Println("URL shortener is running on port 8080")
	 // Run the HTTP server in a separate goroutine
    go func() {
        if err := http.ListenAndServe(":8080", nil); err != nil {
            log.Fatalf("Failed to start server: %v", err)
        }
    }()

	  // Setup signal handling to save URLs on exit
    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt)
    <-c

    if err := urlShortener.SaveToFile("urls.gob"); err != nil {
        fmt.Println("Error saving URLs:", err)
    } else {
        fmt.Println("URLs saved successfully.")
    }
}
