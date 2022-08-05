package main

import (
	"modules/service"

	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {

	r := mux.NewRouter()

	// product := &models.Product{
	// 	Name:        "Product 1",
	// 	Description: "Description 1",
	// 	Price:       1.99,
	// }

	// service.InsertProduct(*product)

	r.HandleFunc("/products", service.GetProducts).Methods("GET")
	r.HandleFunc("/add", service.InsertProduct).Methods("POST")
	r.HandleFunc("/cart/{id}", service.AddToCart).Methods("POST")
	r.HandleFunc("/cart", service.GetCartItems).Methods("GET")
	r.HandleFunc("/deleteCart/{id}", service.DeleteOneItemFromCart).Methods("DELETE")
	r.HandleFunc("/deleteAllCart/{id}", service.DeleteCartItem).Methods("DELETE")

	handler := cors.AllowAll().Handler(r)
	fmt.Printf("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", handler))

}
