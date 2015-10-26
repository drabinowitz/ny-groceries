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
	Id   int64  `json:id`
	Name string `json:name`
}

type Receipt struct {
	id       int64
	store_id int64
	total    float64
	date     time.Time
}

var unitTypes []string = []string{"unit", "qt", "oz", "pt", "lb"}

type Purchase struct {
	id         int64
	receipt_id int64
	quantity   int64
	cost       float64
	product_id int64
	unit       string
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
	http.ListenAndServe(":8080", nil)
}
