package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Product struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

var products = []Product{
	{ID: 1, Name: "Product 1", Price: 10000},
	{ID: 2, Name: "Product 2", Price: 20000},
}

func getProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func createProduct(w http.ResponseWriter, r *http.Request) {
	var p Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	p.ID = getNextID()
	products = append(products, p)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

func updateProduct(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var updated Product
	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	for i, p := range products {
		if p.ID == id {
			updated.ID = id
			products[i] = updated
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updated)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
}

func deleteProduct(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	for i, p := range products {
		if p.ID == id {
			products = append(products[:i], products[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
}

func getNextID() int {
	maxID := 0
	for _, p := range products {
		if p.ID > maxID {
			maxID = p.ID
		}
	}
	return maxID + 1
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/products", getProducts).Methods("GET")
	r.HandleFunc("/products", createProduct).Methods("POST")
	r.HandleFunc("/products/{id}", updateProduct).Methods("PUT")
	r.HandleFunc("/products/{id}", deleteProduct).Methods("DELETE")

	// Wrap router with simple CORS middleware for local development
	handler := corsMiddleware(r)

	log.Println("Backend running at :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}

// corsMiddleware allows requests from localhost:3000 for development (simple, not for production)
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "https://romanza-fullstack-test.vercel.app/")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
