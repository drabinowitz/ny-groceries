package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

type Store struct {
	Id   int64  `json:id`
	Name string `json:name`
}

type Receipt struct {
	Id       int64   `json:id`
	Store_id int64   `json:store_id`
	Total    float64 `json:total`
	Date     string  `json:date`
}

type Purchase struct {
	Id         int64   `json:id`
	Receipt_id int64   `json:receipt_id`
	Quantity   int64   `json:quantity`
	Cost       float64 `json:cost`
	Product_id int64   `json:product_id`
	Unit       string  `json:unit`
}

type Product struct {
	Id           int64  `json:id`
	Category     string `json:category`
	Sub_category string `json:sub_category, omitempty`
}

type ReceiptUpload struct {
	Receipt   Receipt
	Purchases []Purchase
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>STUFF</h1><div>STUFF</div>")
}

var AllStores []Store = make([]Store, 0)

func storesHandler(w http.ResponseWriter, r *http.Request) {
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
			statement := "INSERT INTO products (%s) VALUES (%s) RETURNING id"
			if product.Sub_category != "" {
				data := fmt.Sprintf("'%s', '%s'", product.Category, product.Sub_category)
				statement = fmt.Sprintf(statement, "category, sub_category", data)
			} else {
				data := fmt.Sprintf("'%s'", product.Category)
				statement = fmt.Sprintf(statement, "category", data)
			}
			var productId int64
			err = db.QueryRow(statement).Scan(&productId)
			if err != nil {
				log.Fatal(err)
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

func receiptUploadsHandler(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "PUT" {
			decoder := json.NewDecoder(r.Body)
			var receiptUpload ReceiptUpload
			err := decoder.Decode(&receiptUpload)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			if len(receiptUpload.Purchases) == 0 {
				http.Error(w, "No purchases with receipt", http.StatusInternalServerError)
			}
			receipt := receiptUpload.Receipt
			receiptStatement := fmt.Sprintf(`
				INSERT INTO receipts (store_id, total, date)
				VALUES (%v, %v, '%s') RETURNING id
			`, receipt.Store_id, receipt.Total, receipt.Date)
			var receiptId int64
			err = db.QueryRow(receiptStatement).Scan(&receiptId)
			if err != nil {
				log.Fatal(err)
			}
			receiptUpload.Receipt.Id = receiptId
			purchases := receiptUpload.Purchases
			purchaseString := `
				INSERT INTO purchases (receipt_id, quantity, cost, product_id, unit)
				VALUES (%v, %v, %v, %v, '%s') RETURNING id
			`
			for i, p := range purchases {
				purchaseStatement := fmt.Sprintf(purchaseString,
					receiptId, p.Quantity, p.Cost, p.Product_id, p.Unit)
				var purchaseId int64
				err = db.QueryRow(purchaseStatement).Scan(&purchaseId)
				if err != nil {
					log.Fatal(err)
				}
				receiptUpload.Purchases[i].Id = purchaseId
				receiptUpload.Purchases[i].Receipt_id = receiptId
			}
			js, err := json.Marshal(receiptUpload)
			w.Header().Set("Content-Type", "application/json")
			w.Write(js)
		}
	}
}

func main() {
	db, err := sql.Open("postgres", "postgres:///groceries_test")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	loadStores(db)
	loadProducts(db)

	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/stores/", storesHandler)
	http.HandleFunc("/products/", productsHandler(db))
	http.HandleFunc("/receipt_uploads/", receiptUploadsHandler(db))
	http.ListenAndServe(":8080", nil)
}
