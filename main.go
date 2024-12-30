package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type URL struct {
	id           string    `json:"id"`
	orginalURL   string    `json:"orginalURL"`
	encryptedURL string    `json:"encryptedURL"`
	createdDate  time.Time `json:"createdDate"`
}

var urlDB = make(map[string]URL)

func generateShortURL(orginalURL string) string {
	hasher := md5.New()
	hasher.Write([]byte(orginalURL))
	data := hasher.Sum(nil)
	hash := hex.EncodeToString(data)
	shorterHash := hash[:6]

	return shorterHash
}

func createURL(OrginalURL string) string {
	ShortURL := generateShortURL(OrginalURL)
	Id := ShortURL

	urlDB[Id] = URL{
		id:           Id,
		orginalURL:   OrginalURL,
		encryptedURL: ShortURL,
		createdDate:  time.Now(),
	}
	return ShortURL
}

func getURL(id string) (URL, error) {
	url, ok := urlDB[id]
	if !ok {
		return URL{}, errors.New("URL not found")
	}
	return url, nil
}

func RootPageURL(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello World")
}

func shortURLHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		URL string `json:"url"`
	}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}
	shortURL_ := createURL(data.URL)

	response := struct {
		ShortURL string `json:"short_url"`
	}{ShortURL: shortURL_}

	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(response)
}
func redirectURLHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/redirect/"):]
	url, err := getURL(id)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusNotFound)
	}
	http.Redirect(w, r, url.orginalURL, http.StatusFound)
}

func main() {
	http.HandleFunc("/", RootPageURL)
	http.HandleFunc("/shorten", shortURLHandler)
	http.HandleFunc("/redirect/", redirectURLHandler)

	fmt.Println("Server staritng at 3000....")
	var err error = http.ListenAndServe(":3000", nil)
	if err != nil {
		fmt.Println("The error occured in 'ListenAndServe' function", err)
	}
}
