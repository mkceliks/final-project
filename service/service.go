package service

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"modules/db"
	"modules/models"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetProducts(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query("SELECT * FROM products")
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("No rows found")
			return
		}
		log.Fatal(err)
	}
	defer rows.Close()

	var products []*models.Product
	for rows.Next() {
		prd := &models.Product{}
		err := rows.Scan(&prd.ID, &prd.Name, &prd.Description, &prd.Price, &prd.Discount, &prd.Tax)
		if err != nil {
			log.Fatal(err)
		}
		products = append(products, prd)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(products)
}

func GetProductById(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	row := db.DB.QueryRow("SELECT * FROM products WHERE id = $1", id)
	prd := &models.Product{}
	err := row.Scan(&prd.ID, &prd.Name, &prd.Description, &prd.Price, &prd.Discount, &prd.Tax)
	switch {
	case err == sql.ErrNoRows:
		fmt.Println("No rows found")
		return
	case err != nil:
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(prd)
}

func InsertProduct(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		w.Header().Set("Content-Type", "application/json")
		var prd models.Product
		_ = json.NewDecoder(r.Body).Decode(&prd)
		result, err := db.DB.Exec("INSERT INTO products(name, description, price, discount, tax) VALUES($1, $2, $3, $4, $5)", prd.Name, prd.Description, prd.Price, prd.Discount, prd.Tax)
		if err != nil {
			log.Fatal(err)
		}
		count, err := result.RowsAffected()
		if err != nil {
			log.Fatal(count)
		}
		w.Write([]byte("Product added successfully"))
	}
}

func AddToCart(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		w.Header().Set("Content-Type", "application/json")
		var cart models.Cart
		_ = json.NewDecoder(r.Body).Decode(&cart)
		id := mux.Vars(r)["id"]
		var product string
		err := db.DB.QueryRow("SELECT price FROM products WHERE id = $1", id).Scan(&product) // Get the price of the product
		switch {
		case err == sql.ErrNoRows:
			fmt.Println("No rows found")
			return
		case err != nil:
			log.Fatal(err)
		}
		floatVar, _ := strconv.ParseFloat(product, 64)
		total_price := floatVar * float64(cart.Quantity) //total price of the product

		if db.DB.QueryRow("SELECT * FROM cart WHERE product_id = $1", id).Scan(&cart.ProductID) == sql.ErrNoRows {
			result, err := db.DB.Exec("INSERT INTO cart(product_id, quantity, total_price, total_discount) VALUES($1, $2, $3, $4)", id, cart.Quantity, total_price, cart.Discount)
			if err != nil {
				log.Fatal(err)
			}
			count, err := result.RowsAffected()
			if err != nil {
				log.Fatal(count)
			}
			w.Write([]byte("Product added to cart successfully"))
		} else {
			var oldQuantity int
			err := db.DB.QueryRow("SELECT quantity FROM cart WHERE product_id = $1", id).Scan(&oldQuantity)
			switch {
			case err == sql.ErrNoRows:
				fmt.Println("No rows found")
				return
			case err != nil:
				log.Fatal(err)
			}
			newQuantity := cart.Quantity + oldQuantity
			total_price = floatVar * float64(newQuantity)
			result, err := db.DB.Exec("UPDATE cart SET quantity = $1, total_price = $2, total_discount = $3 WHERE product_id = $4", newQuantity, total_price, cart.Discount, id)
			if err != nil {
				log.Fatal(err)
			}
			count, err := result.RowsAffected()
			if err != nil {
				log.Fatal(count)
			}
			w.Write([]byte("Product added to cart sssssuccessfully"))
		}

	}

}

func AddOneItemToCart(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var cart models.Cart
	_ = json.NewDecoder(r.Body).Decode(&cart)
	id := mux.Vars(r)["id"]
	var product string
	err := db.DB.QueryRow("SELECT price FROM products WHERE id = $1", id).Scan(&product) // Get the price of the product
	switch {
	case err == sql.ErrNoRows:
		fmt.Println("No rows found")
		return
	case err != nil:
		log.Fatal(err)
	}

	if db.DB.QueryRow("SELECT * FROM cart WHERE product_id = $1", id).Scan(&cart.ProductID) != sql.ErrNoRows {
		result, err := db.DB.Exec("UPDATE cart SET quantity = quantity + 1, total_price = total_price + $1 WHERE product_id = $2", product, id)
		if err != nil {
			log.Fatal(err)
		}
		count, err := result.RowsAffected()
		if err != nil {
			log.Fatal(count)
		}
		w.Write([]byte("Cart item added successfully"))
		return
	} else {
		result, err := db.DB.Exec("INSERT INTO cart(product_id, quantity, total_price, total_discount) VALUES($1, $2 + 1 , $3, $4)", id, cart.Quantity, product, cart.Discount)
		if err != nil {
			log.Fatal(err)
		}
		count, err := result.RowsAffected()
		if err != nil {
			log.Fatal(count)
		}
		w.Write([]byte("Cart item added successfully"))
	}

}

func GetCartItems(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query("SELECT * FROM cart")
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("No rows found")
			return
		}
		log.Fatal(err)
	}
	defer rows.Close()

	var cartItems []*models.Cart
	for rows.Next() {
		crt := &models.Cart{}
		err := rows.Scan(&crt.ID, &crt.ProductID, &crt.Quantity, &crt.Price, &crt.Discount)
		if err != nil {
			log.Fatal(err)
		}
		cartItems = append(cartItems, crt)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(cartItems)
}

func DeleteCartItem(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	result, err := db.DB.Exec("DELETE FROM cart WHERE product_id = $1", id)
	if err != nil {
		log.Fatal(err)
	}
	count, err := result.RowsAffected()
	if err != nil {
		log.Fatal(count)
	}
	w.Write([]byte("Cart item deleted successfully"))
}

func DeleteOneItemFromCart(w http.ResponseWriter, r *http.Request) {

	if r.Method == "DELETE" {
		w.Header().Set("Content-Type", "application/json")
		var cart models.Cart
		_ = json.NewDecoder(r.Body).Decode(&cart)
		id := mux.Vars(r)["id"]
		var product string
		err := db.DB.QueryRow("SELECT price FROM products WHERE id = $1", id).Scan(&product) // Get the price of the product
		switch {
		case err == sql.ErrNoRows:
			fmt.Println("No rows found")
			return
		case err != nil:
			log.Fatal(err)
		}

		row := db.DB.QueryRow("SELECT * FROM cart WHERE product_id = $1", id).Scan(&cart.ID, &cart.ProductID, &cart.Quantity, &cart.Price, &cart.Discount)
		switch {
		case row == sql.ErrNoRows:
			fmt.Println("No rows found")
			w.Write([]byte("No Rows Found"))
			return
		case row != nil:
			log.Fatal(row)
		}
		floatVar, _ := strconv.ParseFloat(product, 64)
		total_price := floatVar * (float64(cart.Quantity) - 1) //total price of the product
		newQuantity := cart.Quantity - 1
		if newQuantity <= 0 {
			result, err := db.DB.Exec("DELETE FROM cart WHERE product_id = $1", id)
			if err != nil {
				log.Fatal(err)
			}
			count, err := result.RowsAffected()
			if err != nil {
				log.Fatal(count)
			}
			w.Write([]byte("There is no item more in the cart"))
			return
		} else {
			result, err := db.DB.Exec("UPDATE cart SET quantity = $1, total_price = $2, total_discount = $3 WHERE product_id = $4", newQuantity, total_price, cart.Discount, id)
			if err != nil {
				log.Fatal(err)
			}
			count, err := result.RowsAffected()
			if err != nil {
				log.Fatal(count)
			}
			w.Write([]byte("One item deleted from cart successfully"))
			return
		}
	}
}
