package main

import (
	"encoding/json"
	"github.com/drabinowitz/ny-groceries/api/apidb"
	"net/http"
)

func setHeaders(w http.ResponseWriter, r *http.Request) {
	if origin := r.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
}

func storesHandler(api *apidb.Api) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		setHeaders(w, r)
		js, err := json.Marshal(api.GetAllStores())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}

func productsHandler(api *apidb.Api) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		setHeaders(w, r)
		if r.Method == "GET" {
			js, err := json.Marshal(api.GetAllProducts())
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(js)

		} else if r.Method == "PUT" {
			decoder := json.NewDecoder(r.Body)
			var product apidb.Product
			err := decoder.Decode(&product)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			js, err := json.Marshal(api.AddProduct(product))
			w.Header().Set("Content-Type", "application/json")
			w.Write(js)
		}
	}
}

func receiptUploadsHandler(api *apidb.Api) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		setHeaders(w, r)
		if r.Method == "PUT" {
			decoder := json.NewDecoder(r.Body)
			var receiptUpload apidb.ReceiptUpload
			err := decoder.Decode(&receiptUpload)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			js, err := json.Marshal(api.AddReceiptUpload(receiptUpload))
			w.Header().Set("Content-Type", "application/json")
			w.Write(js)
		}
	}
}

func main() {
	api := apidb.Open()
	defer api.Close()

	http.HandleFunc("/stores/", storesHandler(api))
	http.HandleFunc("/products/", productsHandler(api))
	http.HandleFunc("/receipt_uploads/", receiptUploadsHandler(api))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w,r, "../ui/build/index.html")
    })
	http.ListenAndServe(":8000", nil)
}
