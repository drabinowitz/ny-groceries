package main

import (
	"encoding/json"
	"fmt"
	"github.com/drabinowitz/ny-groceries/api/apidb"
	"net/http"
	"strconv"
	"strings"
)

type UnitCost struct {
	Unit     string  `json:"unit"`
	Quantity float64 `json:"quantity"`
	Cost     float64 `json:"cost"`
}

type StoreCostsByUnit struct {
	Store_id  int64               `json:"store_id"`
	UnitCosts map[string]UnitCost `json:"units"`
}

type ProductCostsByStore struct {
	Product_id int64                       `json:"product_id"`
	StoreCosts map[string]StoreCostsByUnit `json:"stores"`
}

type RequestedProducts struct {
	Products map[string]ProductCostsByStore `json:"products"`
}

func newStoreCostsByUnit(store_id int64) (storeCostsByUnit StoreCostsByUnit) {
	storeCostsByUnit.Store_id = store_id
	storeCostsByUnit.UnitCosts = make(map[string]UnitCost)
	return
}

func newProductCostsByStore(product_id int64) (productCostsByStore ProductCostsByStore) {
	productCostsByStore.Product_id = product_id
	productCostsByStore.StoreCosts = make(map[string]StoreCostsByUnit)
	return
}

func newRequestedProducts() (requestedProducts RequestedProducts) {
	requestedProducts.Products = make(map[string]ProductCostsByStore)
	return
}

func setHeaders(w http.ResponseWriter, r *http.Request) {
	if origin := r.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
}

func productsHandler(api *apidb.Api) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		setHeaders(w, r)
		path := r.URL.Path
		path_arr := strings.Split(path, "/")
		if len(path_arr) == 3 {
			js, err := json.Marshal(api.GetAllProducts())
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(js)
		} else {
			category := strings.ToLower(path_arr[2])
			var sub_category string
			if len(path_arr) == 5 {
				sub_category = strings.ToLower(path_arr[3])
			}
			allProducts := api.GetAllProducts()
			requestedProducts := newRequestedProducts()
			for _, product := range allProducts {
				if strings.ToLower(product.Category) == category {
					if sub_category == "" || strings.ToLower(product.Sub_category) == sub_category {
						requestedProducts.Products[strconv.FormatInt(product.Id, 10)] = newProductCostsByStore(product.Id)
					}
				}
			}
			if len(requestedProducts.Products) == 0 {
				http.Error(w, fmt.Sprintf("found no products matching: '%[1] %[2]'", category, sub_category), 400)
				return
			}
			allPurchases := api.GetAllPurchases()
			allReceipts := api.GetAllReceipts()
			for _, purchase := range allPurchases {
				product_id := strconv.FormatInt(purchase.Product_id, 10)
				productCostsByStore, ok := requestedProducts.Products[product_id]
				if ok {
					var receipt apidb.Receipt
					for _, receiptItr := range allReceipts {
						if receiptItr.Id == purchase.Receipt_id {
							receipt = receiptItr
						}
					}
					store_id := strconv.FormatInt(receipt.Store_id, 10)
					_, ok := productCostsByStore.StoreCosts[store_id]
					if !ok {
						productCostsByStore.StoreCosts[store_id] = newStoreCostsByUnit(receipt.Store_id)
					}
					unitCost := productCostsByStore.StoreCosts[store_id].UnitCosts[purchase.Unit]
					unitCost.Unit = purchase.Unit
					unitCost.Quantity += purchase.Quantity
					unitCost.Cost += purchase.Cost
					requestedProducts.Products[product_id].StoreCosts[store_id].UnitCosts[purchase.Unit] = unitCost
				}
			}
			js, err := json.Marshal(requestedProducts)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(js)
		}
	}
}

func main() {
	api := apidb.Open()
	defer api.Close()

	http.HandleFunc("/products/", productsHandler(api))
	http.ListenAndServe(":8000", nil)
}
