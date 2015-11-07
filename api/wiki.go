package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"time"
)

type Store struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type Receipt struct {
	Id       int64   `json:"id"`
	Store_id int64   `json:"store_id"`
	Total    float64 `json:"total"`
	Date     string  `json:"date"`
}

type Purchase struct {
	Id         int64   `json:"id"`
	Receipt_id int64   `json:"receipt_id"`
	Quantity   float64 `json:"quantity"`
	Cost       float64 `json:"cost"`
	Product_id int64   `json:"product_id"`
	Unit       string  `json:"unit"`
}

type Product struct {
	Id           int64  `json:"id"`
	Category     string `json:"category"`
	Sub_category string `json:"sub_category,omitempty"`
}

type ReceiptUpload struct {
	Receipt   Receipt    `json:"receipt"`
	Purchases []Purchase `json:"purchases"`
}

var AllStores []Store = make([]Store, 0)

func storesHandler(w http.ResponseWriter, r *http.Request) {
	if origin := r.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	js, err := json.Marshal(AllStores)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func loadStores(db *sql.DB) {
	stores, err := db.Query("SELECT id, name FROM stores")
	if err != nil {
		log.Fatal(err)
	}
	defer stores.Close()
	var (
		id   int64
		name string
	)
	for stores.Next() {
		err := stores.Scan(&id, &name)
		if err != nil {
			log.Fatal(err)
		}
		AllStores = append(AllStores, Store{id, name})
	}
}

var AllProducts []Product = make([]Product, 0)

func productsHandler(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		if r.Method == "GET" {
			js, err := json.Marshal(AllProducts)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(js)

		} else if r.Method == "PUT" {
			decoder := json.NewDecoder(r.Body)
			var product Product
			err := decoder.Decode(&product)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			if product.Category == "" {
				http.Error(w, "Category Required", http.StatusInternalServerError)
			}
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			statement := "INSERT INTO products (%s) VALUES (%s)"
			if product.Sub_category != "" {
				data := fmt.Sprintf("'%s', '%s'", product.Category, product.Sub_category)
				statement = fmt.Sprintf(statement, "category, sub_category", data)
			} else {
				data := fmt.Sprintf("'%s'", product.Category)
				statement = fmt.Sprintf(statement, "category", data)
			}
			var productId int64
			_, err = db.Exec(statement)
			productIds, err := db.Query("SELECT last_insert_rowid() FROM products")
			if err != nil {
				log.Fatal(err)
			}
			defer productIds.Close()
			for productIds.Next() {
				err := productIds.Scan(&productId)
				if err != nil {
					log.Fatal(err)
				}
			}
			product.Id = productId
			AllProducts = append(AllProducts, product)
			js, err := json.Marshal(product)
			w.Header().Set("Content-Type", "application/json")
			w.Write(js)
		}
	}
}

func loadProducts(db *sql.DB) {
	products, err := db.Query("SELECT id, category, sub_category FROM products")
	if err != nil {
		log.Fatal(err)
	}
	defer products.Close()
	var (
		id           int64
		category     string
		sub_category sql.NullString
	)
	for products.Next() {
		err := products.Scan(&id, &category, &sub_category)
		if err != nil {
			log.Fatal(err)
		}
		product := Product{Id: id, Category: category}
		if sub_category.Valid {
			product.Sub_category = sub_category.String
		}
		AllProducts = append(AllProducts, product)
	}
}

func newPurchase() *Purchase {
	return &Purchase{
		Quantity: 1,
		Unit:     "unit",
	}
}

func receiptUploadsHandler(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		if r.Method == "PUT" {
			decoder := json.NewDecoder(r.Body)
			var receiptUpload ReceiptUpload
			err := decoder.Decode(&receiptUpload)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			receipt := receiptUpload.Receipt
			receiptDate, err := time.Parse("01/02/2006", receipt.Date)
			if err != nil {
				log.Fatal(err)
			}
			receiptStatement := fmt.Sprintf(`
				INSERT INTO receipts (store_id, total, date)
				VALUES (%v, %v, '%s')
			`, receipt.Store_id, receipt.Total, receiptDate.Format("2006-01-02"))
			var receiptId int64
			_, err = db.Exec(receiptStatement)
			if err != nil {
				log.Fatal(err)
			}
			receiptIds, err := db.Query("SELECT last_insert_rowid() FROM receipts")
			if err != nil {
				log.Fatal(err)
			}
			defer receiptIds.Close()
			for receiptIds.Next() {
				err := receiptIds.Scan(&receiptId)
				if err != nil {
					log.Fatal(err)
				}
			}
			receiptUpload.Receipt.Id = receiptId
			purchases := receiptUpload.Purchases
			purchaseString := `
				INSERT INTO purchases (receipt_id, quantity, cost, product_id, unit)
				VALUES (%v, %v, %v, %v, '%s')
			`
			for i := range purchases {
				p := &receiptUpload.Purchases[i]
				if p.Quantity == 0 {
					p.Quantity = 1
				}
				if p.Unit == "" {
					p.Unit = "unit"
				}
				purchaseStatement := fmt.Sprintf(purchaseString,
					receiptId, p.Quantity, p.Cost, p.Product_id, p.Unit)
				var purchaseId int64
				_, err = db.Exec(purchaseStatement)
				if err != nil {
					log.Fatal(err)
				}
				purchaseIds, err := db.Query("SELECT last_insert_rowid() FROM purchases")
				if err != nil {
					log.Fatal(err)
				}
				defer purchaseIds.Close()
				for purchaseIds.Next() {
					err := purchaseIds.Scan(&purchaseId)
					if err != nil {
						log.Fatal(err)
					}
				}
				p.Id = purchaseId
				p.Receipt_id = receiptId
			}
			js, err := json.Marshal(receiptUpload)
			w.Header().Set("Content-Type", "application/json")
			w.Write(js)
		}
	}
}

func main() {
	db, err := sql.Open("sqlite3", "./apidb.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	loadStores(db)
	loadProducts(db)

	http.HandleFunc("/stores/", storesHandler)
	http.HandleFunc("/products/", productsHandler(db))
	http.HandleFunc("/receipt_uploads/", receiptUploadsHandler(db))
	http.ListenAndServe(":8000", nil)
}
