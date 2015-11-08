package main

import (
	"./apidb.go"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"time"
)

func setHeaders(w *http.ResponseWriter) {
	if origin := r.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
}

func storesHandler(db *sql.db) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		setHeaders(&w)
		js, err := json.Marshal(db.GetAllStores())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}

func productsHandler(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		setHeaders(&w)
		if r.Method == "GET" {
			js := json.Marshal(db.GetAllProducts())
			w.Header().Set("Content-Type", "application/json")
			w.Write(js)

		} else if r.Method == "PUT" {
			decoder := json.NewDecoder(r.Body)
			var product apidb.Product
			err := decoder.Decode(&product)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			js, err := json.Marshal(db.AddProduct(product))
			w.Header().Set("Content-Type", "application/json")
			w.Write(js)
		}
	}
}

func receiptUploadsHandler(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		setHeaders(&w)
		if r.Method == "PUT" {
			decoder := json.NewDecoder(r.Body)
			var receiptUpload apidb.ReceiptUpload
			err := decoder.Decode(&receiptUpload)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			js, err := json.Marshal(db.AddReceiptUpload(receiptUpload))
			w.Header().Set("Content-Type", "application/json")
			w.Write(js)
		}
	}
}

func main() {
	db := apidb.Open()
	defer db.Close()

	http.HandleFunc("/stores/", storesHandler(db))
	http.HandleFunc("/products/", productsHandler(db))
	http.HandleFunc("/receipt_uploads/", receiptUploadsHandler(db))
	http.ListenAndServe(":8000", nil)
}
