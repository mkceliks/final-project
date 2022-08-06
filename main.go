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

	r.HandleFunc("/products", service.GetProducts).Methods("GET")               // GET all products
	r.HandleFunc("/cart", service.GetCartItems).Methods("GET")                  // GET cart items
	r.HandleFunc("/getProductById/{id}", service.GetProductById).Methods("GET") // GET product by id

	r.HandleFunc("/addProduct", service.InsertProduct).Methods("POST")               // POST a new product to the database. Requires a JSON body with the product details ( Quantity ).
	r.HandleFunc("/addToCart/{id}", service.AddToCart).Methods("POST")               // add to cart
	r.HandleFunc("/addOneItemToCart/{id}", service.AddOneItemToCart).Methods("POST") // POST one item to cart

	r.HandleFunc("/deleteOneItemFromCart/{id}", service.DeleteOneItemFromCart).Methods("DELETE") // DELETE one item from cart
	r.HandleFunc("/deleteAllCart/{id}", service.DeleteCartItem).Methods("DELETE")                // DELETE all cart items

	handler := cors.AllowAll().Handler(r)
	fmt.Printf("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", handler))

}
