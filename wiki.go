package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"time"
)

type Store struct {
	Id   int32  `json:id`
	Name string `json:name`
}

type Receipt struct {
	id       int32
	store_id int32
	total    float32
	date     time.Time
}

var unitTypes []string = []string{"unit", "qt", "oz", "pt", "lb"}

type Purchase struct {
	id         int32
	receipt_id int32
	quantity   int32
	cost       float32
	product_id int32
	unit       string
}

type Product struct {
	Id           int32  `json:id`
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
		id   int32
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

func productsHandler(w http.ResponseWriter, r *http.Request) {
	js, err := json.Marshal(AllProducts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func loadProducts(db *sql.DB) {
	products, err := db.Query("SELECT id, category, sub_category FROM products")
	if err != nil {
		log.Fatal(err)
	}
	defer products.Close()
	var (
		id           int32
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

func main() {
	db, err := sql.Open("postgres", "postgres:///groceries")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	loadStores(db)
	loadProducts(db)

	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/stores/", storesHandler)
	http.HandleFunc("/products/", productsHandler)
	http.ListenAndServe(":8080", nil)
}
