package main

import (
	"encoding/json"
	"github.com/drabinowitz/ny-groceries/api/apidb"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
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

func purchasesHandler(api *apidb.Api) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		setHeaders(w, r)
		vars := mux.Vars(r)
		productId, err := strconv.ParseInt(vars["productId"], 10, 64)
		storeId, err := strconv.ParseInt(vars["storeId"], 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		stores := api.GetAllStores()
		var foundStore apidb.Store
		for _, store := range stores {
			if store.Id == storeId {
				foundStore = store
				break
			}
		}

		products := api.GetAllProducts()
		var foundProduct apidb.Product
		for _, product := range products {
			if product.Id == productId {
				foundProduct = product
				break
			}
		}

		if foundStore.Id != storeId || foundProduct.Id != productId {
			http.Error(w, "cannot find store or product", http.StatusInternalServerError)
			return
		}

		receipts := api.GetAllReceipts()
		var foundReceipts []apidb.Receipt = make([]apidb.Receipt, 0)
		for _, receipt := range receipts {
			if receipt.Store_id == foundStore.Id {
				foundReceipts = append(foundReceipts, receipt)
			}
		}

		purchases := api.GetAllPurchases()
		var foundPurchases []apidb.Purchase = make([]apidb.Purchase, 0)
		for _, purchase := range purchases {
			if purchase.Product_id == foundProduct.Id {
				for _, receipt := range foundReceipts {
					if purchase.Receipt_id == receipt.Id {
						foundPurchases = append(foundPurchases, purchase)
						break
					}
				}
			}
		}

		w.Header().Set("Content-Type", "application/json")
		js, err := json.Marshal(foundPurchases)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
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

	r := mux.NewRouter()
	r.HandleFunc("/stores/", storesHandler(api))
	r.HandleFunc("/purchases/{storeId}/{productId}", purchasesHandler(api))
	r.HandleFunc("/products/", productsHandler(api))
	r.HandleFunc("/receipt_uploads/", receiptUploadsHandler(api))
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../ui/build/index.html")
	})
	http.ListenAndServe(":8000", r)
}
