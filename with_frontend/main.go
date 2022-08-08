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
	r.HandleFunc("/customers", service.GetCustomers).Methods("GET")             // GET all customers
	r.HandleFunc("/getProductById/{id}", service.GetProductById).Methods("GET") // GET product by id (id = productid(int))

	r.HandleFunc("/addProduct", service.InsertProduct).Methods("POST")   // POST a new product to the database. Requires a JSON body with the product details ( Quantity ).
	r.HandleFunc("/addCustomer", service.InsertCustomer).Methods("POST") // POST a new customer to the database. Requires a JSON body with the customer details ( UserName ).

	r.HandleFunc("/addToCart/{id}", service.AddToCart).Methods("POST")               // add to cart by id. Requires a JSON body with the product details ( Quantity ).
	r.HandleFunc("/addOneItemToCart/{id}", service.AddOneItemToCart).Methods("POST") // POST one item to cart ( id = product id(int) )

	r.HandleFunc("/deleteOneItemFromCart/{id}", service.DeleteOneItemFromCart).Methods("DELETE") // DELETE one item from cart ( id = productid(int) )
	r.HandleFunc("/deleteAllCart/{id}", service.DeleteCartItem).Methods("DELETE")                // DELETE the row from cart ( id =  productid(int))

	r.HandleFunc("/order", service.Order).Methods("POST") // POST purchase order ( JSON: CustomerID: int )

	handler := cors.AllowAll().Handler(r)
	fmt.Printf("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", handler))

}
