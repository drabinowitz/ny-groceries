package main

import (
	"encoding/json"
	"fmt"
	"github.com/drabinowitz/ny-groceries/api/apidb"
	"net/http"
	"sort"
	"strings"
)

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
		if len(path_arr) == 2 {
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
			if len(path_arr) == 4 {
				sub_category = strings.ToLower(path_arr[3])
			}
			allProducts := api.GetAllProducts()
			requestedProducts := make(map[int64]map[int64]map[string][]float64)
			for _, product := range allProducts {
				if strings.ToLower(product.Category) == category {
					if sub_category == "" || strings.ToLower(product.Sub_category) == sub_category {
						requestedProducts[product.Id] = make(map[int64]map[string][]float64)
					}
				}
			}
			if len(requestedProducts) == 0 {
				http.Error(w, fmt.Sprintf("found no products matching: '%[1] %[2]'", category, sub_category), 400)
				return
			}
			allPurchases := api.GetAllPurchases()
			allReceipts := api.GetAllReceipts()
			for _, purchase := range allPurchases {
				product_id := purchase.Product_id
				unitCostsByStore, ok := requestedProducts[product_id]
				if ok {
					idx := sort.Search(len(allReceipts), func(i int) bool {
						return allReceipts[i].Id == purchase.Receipt_id
					})
					receipt := allReceipts[idx]
					unitCost := unitCostsByStore[receipt.Store_id][purchase.Unit]
					unitCost = append(unitCost, float64(purchase.Cost)/float64(purchase.Quantity))
					requestedProducts[product_id][receipt.Store_id][purchase.Unit] = unitCost
				}
			}
			costedProducts := make(map[int64]map[int64]map[string]float64)
			for product_id, unitCostsByStore := range requestedProducts {
				for store_id, unitCosts := range unitCostsByStore {
					for unit, costs := range unitCosts {
						var averagedCost float64 = 0
						for _, cost := range costs {
							averagedCost += cost
						}
						averagedCost = averagedCost / float64(len(costs))
						costedProducts[product_id][store_id][unit] = averagedCost
					}
				}
			}
			js, err := json.Marshal(costedProducts)
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
