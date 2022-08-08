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

func GetCustomers(w http.ResponseWriter, r *http.Request) { // Get all Customers
	rows, err := db.DB.Query("SELECT * FROM customer") // SQL Query to get all customers

	//err handlings**********************
	if err != nil {
		log.Fatal(err)
	}
	//err handlings**********************

	defer rows.Close() // Close the rows after we finish using them

	var customers []*models.Customer // Create a slice of customers to store all the customers we get from the database
	for rows.Next() {                // Loop through all the customers in the database
		var customer models.Customer                       // Create a new customer
		err := rows.Scan(&customer.ID, &customer.UserName) // Scan the customer into the customer struct

		//err handlings**********************
		if err != nil {
			log.Fatal(err)
		}
		//err handlings**********************

		customers = append(customers, &customer) // Append the customer to the slice of customers
	}

	//err handlings**********************
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
	//err handlings**********************

	json.NewEncoder(w).Encode(customers) // Encode the slice of customers into JSON and send it to the client
}

func GetProducts(w http.ResponseWriter, r *http.Request) { // Get all products
	rows, err := db.DB.Query("SELECT * FROM products") // SQL Query to get all products

	//err handlings**********************
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("No rows found")
			return
		}
		log.Fatal(err)
	}
	//err handlings**********************

	defer rows.Close() // Close the rows after we finish using them

	var products []*models.Product // Create a slice of products to store all the products we get from the database

	for rows.Next() { // Loop through all the products in the database
		prd := &models.Product{}                                                     // Create a new product
		err := rows.Scan(&prd.ID, &prd.Name, &prd.Description, &prd.Price, &prd.Tax) // Scan the product into the product struct

		//err handlings**********************
		if err != nil {
			log.Fatal(err)
		}
		//err handlings**********************

		products = append(products, prd) // Append the product to the slice of products
	}

	//err handlings**********************
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
	//err handlings**********************

	json.NewEncoder(w).Encode(products) // Encode the slice of products into JSON and send it to the client
}

func GetProductById(w http.ResponseWriter, r *http.Request) { // Get a product by id
	id := mux.Vars(r)["id"]                                                     // Get the id from the url
	row := db.DB.QueryRow("SELECT * FROM products WHERE id = $1", id)           // SQL Query to get the product by id
	prd := &models.Product{}                                                    // Create a new product
	err := row.Scan(&prd.ID, &prd.Name, &prd.Description, &prd.Price, &prd.Tax) // Scan the product into the product struct

	//err handlings**********************
	switch {
	case err == sql.ErrNoRows:
		fmt.Println("No rows found")
		return
	case err != nil:
		log.Fatal(err)
	}
	//err handlings**********************

	json.NewEncoder(w).Encode(prd) // Encode the product into JSON and send it to the client
}

func InsertProduct(w http.ResponseWriter, r *http.Request) { // Insert a product to the database

	if r.Method == "POST" {
		w.Header().Set("Content-Type", "application/json")
		var prd models.Product
		_ = json.NewDecoder(r.Body).Decode(&prd) // Decode the JSON into the product struct

		result, err := db.DB.Exec("INSERT INTO products(name, description, price, tax) VALUES($1, $2, $3, $4)", prd.Name, prd.Description, prd.Price, prd.Tax) // Execute the SQL Query to insert the product into the database

		//err handlings**********************
		if err != nil {
			log.Fatal(err)
		}
		count, err := result.RowsAffected()
		if err != nil {
			log.Fatal(count)
		}
		//err handlings**********************

		w.Write([]byte("Product added successfully")) // Send a message to the client that the product was added successfully

	}
}

func InsertCustomer(w http.ResponseWriter, r *http.Request) { // Insert a customer to the database

	if r.Method == "POST" {
		w.Header().Set("Content-Type", "application/json")
		var customer models.Customer
		_ = json.NewDecoder(r.Body).Decode(&customer) // Decode the JSON into the customer struct

		result, err := db.DB.Exec("INSERT INTO customer(username) VALUES($1)", customer.UserName) // Execute the SQL Query to insert the customer into the database

		//err handlings**********************
		if err != nil {
			log.Fatal(err)
		}
		count, err := result.RowsAffected()
		if err != nil {
			log.Fatal(count)
		}
		//err handlings**********************

		w.Write([]byte("Customer added successfully")) // Send a message to the client that the customer was added successfully

	}
}

func AddToCart(w http.ResponseWriter, r *http.Request) { // Add a product to the cart WITH QUANTITY

	if r.Method == "POST" {
		w.Header().Set("Content-Type", "application/json")
		var cart models.Cart
		_ = json.NewDecoder(r.Body).Decode(&cart) // Decode the JSON into the cart struct (QUANTITY)
		id := mux.Vars(r)["id"]                   // Get the id from the url

		var product string
		err := db.DB.QueryRow("SELECT price FROM products WHERE id = $1", id).Scan(&product) // SQL Query to get the original price of the product by id

		//err handlings**********************
		switch {
		case err == sql.ErrNoRows:
			fmt.Println("No rows found")
			return
		case err != nil:
			log.Fatal(err)
		}
		//err handlings**********************

		var tax float64
		err = db.DB.QueryRow("SELECT tax FROM products WHERE id = $1", id).Scan(&tax) // SQL Query to get the tax of the product by id

		//err handlings**********************
		switch {
		case err == sql.ErrNoRows:
			fmt.Println("No rows found")
			return
		case err != nil:
			log.Fatal(err)
		}
		//err handlings**********************

		if db.DB.QueryRow("SELECT * FROM cart WHERE product_id = $1", id).Scan(&cart.ProductID) == sql.ErrNoRows && cart.Quantity > 3 { // Check if the product is already in the cart and the user wants to add more than 3 products
			floatVar, _ := strconv.ParseFloat(product, 64)                                                                                                                                                                          // Convert the string to a float64
			result, err := db.DB.Exec("INSERT INTO cart(product_id, quantity, total_price, total_discount) VALUES($1, $2, $3, $4)", id, cart.Quantity, totalPrice(cart.Quantity, floatVar), totalDiscount(cart.Quantity, floatVar)) // Execute the SQL Query to insert the product into the cart with the quantity and the total price and discount

			//err handlings**********************
			if err != nil {
				log.Fatal(err)
			}
			count, err := result.RowsAffected()
			if err != nil {
				log.Fatal(count)
			}
			//err handlings**********************

			w.Write([]byte("There was no product in the cart, and more than 3 products were added")) // Send a message to the client that the product was added successfully

		} else if db.DB.QueryRow("SELECT * FROM cart WHERE product_id = $1", id).Scan(&cart.ProductID) == sql.ErrNoRows && cart.Quantity <= 3 { // Check if the product is already in the cart and the user wants to add less than 3 products

			floatVar, _ := strconv.ParseFloat(product, 64) // Convert the string to a float64

			result, err := db.DB.Exec("INSERT INTO cart(product_id, quantity, total_price, total_discount) VALUES($1, $2, $3, $4)", id, cart.Quantity, totalPrice(cart.Quantity, floatVar), totalDiscount(cart.Quantity, floatVar)) // Execute the SQL Query to insert the product into the cart with the quantity and the total price and discount

			//err handlings**********************
			if err != nil {
				log.Fatal(err)
			}
			count, err := result.RowsAffected()
			if err != nil {
				log.Fatal(count)
			}
			//err handlings**********************

			w.Write([]byte("There was no product in the cart, and less than 3 products were added")) // Send a message to the client that the product was added successfully

		} else { // There is already a product in the cart

			var oldQuantity int
			err := db.DB.QueryRow("SELECT quantity FROM cart WHERE product_id = $1", id).Scan(&oldQuantity) // SQL Query to get how many products already in the cart by id

			//err handlings**********************
			switch {
			case err == sql.ErrNoRows:
				fmt.Println("No rows found")
				return
			case err != nil:
				log.Fatal(err)
			}
			//err handlings**********************

			floatVar, _ := strconv.ParseFloat(product, 64)                                                                                                                                                                    // Convert the string to a float64 ( product --> original price )
			newQuantity := cart.Quantity + oldQuantity                                                                                                                                                                        // Add the old quantity to the new quantity
			result, err := db.DB.Exec("UPDATE cart SET quantity = $1, total_price = $2, total_discount = $3 WHERE product_id = $4", newQuantity, totalPrice(newQuantity, floatVar), totalDiscount(newQuantity, floatVar), id) // Execute the SQL Query to update the cart with the new quantity and the total price and discount

			//err handlings**********************
			if err != nil {
				log.Fatal(err)
			}
			count, err := result.RowsAffected()
			if err != nil {
				log.Fatal(count)
			}
			//err handlings**********************

			w.Write([]byte("There was some products in the cart, and cart is updated with the new quantity")) // Send a message to the client that the product was added successfully
		}
	}
}
func AddOneItemToCart(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var cart models.Cart
	_ = json.NewDecoder(r.Body).Decode(&cart)
	id := mux.Vars(r)["id"] // Get the id from the url
	var product string
	err := db.DB.QueryRow("SELECT price FROM products WHERE id = $1", id).Scan(&product) // SQL Query to get the original price of the product by id0

	//err handlings**********************
	switch {
	case err == sql.ErrNoRows:
		fmt.Println("No rows found")
		return
	case err != nil:
		log.Fatal(err)
	}
	//err handlings**********************

	var tax float64
	err = db.DB.QueryRow("SELECT tax FROM products WHERE id = $1", id).Scan(&tax) // SQL Query to get the tax of the product by id

	//err handlings**********************
	switch {
	case err == sql.ErrNoRows:
		fmt.Println("No rows found")
		return
	case err != nil:
		log.Fatal(err)
	}
	//err handlings**********************

	if db.DB.QueryRow("SELECT * FROM cart WHERE product_id = $1", id).Scan(&cart.ProductID) == sql.ErrNoRows { // Check if the product is already in the cart, if not
		floatVar, _ := strconv.ParseFloat(product, 64)                                                                                                                                                    // Convert the string to a float64
		result, err := db.DB.Exec("INSERT INTO cart(product_id, quantity, total_price, total_discount) VALUES($1, $2 + 1 , $3, $4)", id, cart.Quantity, floatVar, totalDiscount(cart.Quantity, floatVar)) // Execute the SQL Query to insert the product into the cart with the quantity and the total price and discount

		//err handlings**********************
		if err != nil {
			log.Fatal(err)
		}
		count, err := result.RowsAffected()
		if err != nil {
			log.Fatal(count)
		}
		//err handlings**********************

		w.Write([]byte("There was no item in the cart, the product is added")) // Send a message to the client that the product was added successfully
	} else { // There is already a product in the cart
		err := db.DB.QueryRow("SELECT quantity FROM cart WHERE product_id = $1", id).Scan(&cart.Quantity) // SQL Query to get how many products already in the cart

		//err handlings**********************
		switch {
		case err == sql.ErrNoRows:
			fmt.Println("No rows found")
			return
		case err != nil:
			log.Fatal(err)
		}
		//err handlings**********************

		if db.DB.QueryRow("SELECT * FROM cart WHERE product_id = $1", id).Scan(&cart.ProductID) != sql.ErrNoRows && cart.Quantity >= 3 { // Check if the product is already in the cart and if the quantity is greater than 3
			floatVar, _ := strconv.ParseFloat(product, 64)                                                                                                                                                                 // Convert the string to a float64
			newQuantity := cart.Quantity + 1                                                                                                                                                                               // Add 1 to the quantity
			result, err := db.DB.Exec("UPDATE cart SET quantity = quantity + 1, total_price = $1, total_discount = $2 WHERE product_id = $3", totalPrice(newQuantity, floatVar), totalDiscount(newQuantity, floatVar), id) // Execute the SQL Query to update the cart with the new quantity and the total price and discount

			//err handlings**********************
			if err != nil {
				log.Fatal(err)
			}
			count, err := result.RowsAffected()
			if err != nil {
				log.Fatal(count)
			}
			//err handlings**********************

			w.Write([]byte("There was items in the cart more than 3 or equal to 3, and the quantity is updated by 1")) // Send a message to the client that the product was added successfully
			return
		} else if db.DB.QueryRow("SELECT * FROM cart WHERE product_id = $1", id).Scan(&cart.ProductID) != sql.ErrNoRows && cart.Quantity < 3 { // Check if the product is already in the cart and if the quantity is less than 3
			floatVar, _ := strconv.ParseFloat(product, 64)                                                    // Convert the string to a float64
			err := db.DB.QueryRow("SELECT quantity FROM cart WHERE product_id = $1", id).Scan(&cart.Quantity) // SQL Query to get how many products already in the cart

			//err handlings**********************
			switch {
			case err == sql.ErrNoRows:
				fmt.Println("No rows found")
				return
			case err != nil:
				log.Fatal(err)
			}
			//err handlings**********************

			newQuantity := cart.Quantity + 1                                                                                                                                                                               // Add 1 to the quantity
			result, err := db.DB.Exec("UPDATE cart SET quantity = quantity + 1, total_price = $1, total_discount = $2 WHERE product_id = $3", totalPrice(newQuantity, floatVar), totalDiscount(newQuantity, floatVar), id) // Execute the SQL Query to update the cart with the new quantity and the total price and discount

			//err handlings**********************
			if err != nil {
				log.Fatal(err)
			}
			count, err := result.RowsAffected()
			if err != nil {
				log.Fatal(count)
			}
			//err handlings**********************

			w.Write([]byte("There was items below 3 in the cart, and the quantity is updated by 1")) // Send a message to the client that the product was added successfully
		}
	}
}

func GetCartItems(w http.ResponseWriter, r *http.Request) { // Get all the items in the cart
	rows, err := db.DB.Query("SELECT * FROM cart") // SQL Query to get all the items in the cart

	//err handlings**********************
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("No rows found")
			return
		}
		log.Fatal(err)
	}
	//err handlings**********************

	defer rows.Close() // Close the rows after the query is finished

	var cartItems []*models.Cart // Create a slice of cart items
	for rows.Next() {            // Loop through the rows
		crt := &models.Cart{}                                                               // Create a cart item
		err := rows.Scan(&crt.ID, &crt.ProductID, &crt.Quantity, &crt.Price, &crt.Discount) // Scan the row into the cart item

		//err handlings**********************
		if err != nil {
			log.Fatal(err)
		}
		//err handlings**********************

		cartItems = append(cartItems, crt) // Append the cart item to the slice
	}

	//err handlings**********************
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
	//err handlings**********************

	json.NewEncoder(w).Encode(cartItems) // Encode the cart items into json and send it to the client
}

func DeleteCartItem(w http.ResponseWriter, r *http.Request) { // Delete an item from the cart
	id := mux.Vars(r)["id"]                                                 // Get the id from the url
	result, err := db.DB.Exec("DELETE FROM cart WHERE product_id = $1", id) // Execute the SQL Query to delete the item from the cart

	//err handlings**********************
	if err != nil {
		log.Fatal(err)
	}
	count, err := result.RowsAffected()
	if err != nil {
		log.Fatal(count)
	}
	//err handlings**********************

	w.Write([]byte("Cart item's row deleted successfully")) // Send a message to the client that the cart item was deleted successfully
}

func DeleteOneItemFromCart(w http.ResponseWriter, r *http.Request) { // Delete one item from cart

	if r.Method == "DELETE" {
		w.Header().Set("Content-Type", "application/json")
		var cart models.Cart
		_ = json.NewDecoder(r.Body).Decode(&cart)
		id := mux.Vars(r)["id"] // Get the id of the product from the url
		var product string
		err := db.DB.QueryRow("SELECT price FROM products WHERE id = $1", id).Scan(&product) // Get the price of the product from the database

		//err handlings**********************
		switch {
		case err == sql.ErrNoRows:
			fmt.Println("No rows found")
			return
		case err != nil:
			log.Fatal(err)
		}
		//err handlings**********************

		row := db.DB.QueryRow("SELECT * FROM cart WHERE product_id = $1", id).Scan(&cart.ID, &cart.ProductID, &cart.Quantity, &cart.Price, &cart.Discount) // Get the cart details of the product from the database

		//err handlings**********************
		switch {
		case row == sql.ErrNoRows:
			fmt.Println("No rows found")
			w.Write([]byte("No Rows Found"))
			return
		case row != nil:
			log.Fatal(row)
		}
		//err handlings**********************

		floatVar, _ := strconv.ParseFloat(product, 64) // Convert the price of the product to float64
		newQuantity := cart.Quantity - 1               // Subtract 1 from the quantity of the product

		if newQuantity <= 0 { // If the newQuantity is 0 or less than 0, delete the product from the cart
			result, err := db.DB.Exec("DELETE FROM cart WHERE product_id = $1", id) // Delete the product from the cart SQL query

			//err handlings**********************
			if err != nil {
				log.Fatal(err)
			}
			count, err := result.RowsAffected()
			if err != nil {
				log.Fatal(count)
			}
			//err handlings**********************

			w.Write([]byte("There is no item more in the cart")) // POSTMAN response
			return
		} else if newQuantity > 3 { // If the newQuantity is greater than 3, then the quantity is decreased by 1. total_price and total_discount are calculated again.
			//total price of the product
			result, err := db.DB.Exec("UPDATE cart SET quantity = $1, total_price = $2, total_discount = $3 WHERE product_id = $4", newQuantity, totalPrice(newQuantity, floatVar), totalDiscount(newQuantity, floatVar), id) // Update the cart SQL query --> totalDiscount and totalPrice funcs are used to calculate the new total_price and total_discount of the product

			//err handlings**********************
			if err != nil {
				log.Fatal(err)
			}
			count, err := result.RowsAffected()
			if err != nil {
				log.Fatal(count)
			}
			//err handlings**********************

			w.Write([]byte("One item deleted from the cart and discount updated..")) // POSTMAN response
			return
		} else { // If the newQuantity is less than 3 or equal to 3, then the quantity is decreased by 1. total_price and total_discount are calculated again.
			result, err := db.DB.Exec("UPDATE cart SET quantity = $1, total_price = $2, total_discount = $3 WHERE product_id = $4", newQuantity, totalPrice(newQuantity, floatVar), totalDiscount(newQuantity, floatVar), id) // Update the cart SQL query --> totalDiscount and totalPrice funcs are used to calculate the new total_price and total_discount of the product

			//err handlings**********************
			if err != nil {
				log.Fatal(err)
			}
			count, err := result.RowsAffected()
			if err != nil {
				log.Fatal(count)
			}
			//err handlings**********************

			w.Write([]byte("One item deleted from cart successfully")) // POSTMAN response
			return
		}
	}
}

func Order(w http.ResponseWriter, r *http.Request) { // Completing the shopping function

	if r.Method == "POST" {
		w.Header().Set("Content-Type", "application/json")
		var order_details models.OrderDetails
		_ = json.NewDecoder(r.Body).Decode(&order_details) // Decode the json data from the client

		rows, err := db.DB.Query("SELECT * FROM order_details WHERE customer_id = $1", order_details.CustomerID) // Get the order details of the customer from the database

		//err handlings**********************
		if err != nil {
			if err == sql.ErrNoRows {
				fmt.Println("No rows found")
				return
			}
			log.Fatal(err)
		}
		//err handlings**********************

		defer rows.Close() // Close the rows after the query is finished

		var curr_ordertable []*models.OrderDetails // Create a slice of order_details structs to store the order details of the customer
		for rows.Next() {                          // Loop through the rows and add the order details of the customer to the curr_ordertable slice
			ord := &models.OrderDetails{}                               // Create a new order_details struct
			err := rows.Scan(&ord.ID, &ord.CustomerID, &ord.TotalPrice) // Scan the order details of the customer to the order_details struct

			//err handlings**********************
			if err != nil {
				log.Fatal(err)
			}
			//err handlings**********************

			curr_ordertable = append(curr_ordertable, ord) // Add the order details of the customer to the curr_ordertable slice
		}

		//err handlings**********************
		if err = rows.Err(); err != nil {
			log.Fatal(err)
		}
		//err handlings**********************

		var count_Arr []float32
		for _, element := range curr_ordertable { // Loop through the curr_ordertable slice and check if the customer has already placed an order more than price of 100
			if element.TotalPrice > 100 { // given amount of money
				count_Arr = append(count_Arr, element.TotalPrice)
			}
		}

		if len(count_Arr) < 1 { // If the length of the count_Arr is less than 1, than the customer didn't place an order more than 100 so no discount is applied
			rows, err := db.DB.Query("SELECT * FROM cart") // Get the cart of the customer from the database

			//err handlings**********************
			if err != nil {
				if err == sql.ErrNoRows {
					fmt.Println("No rows found")
					return
				}
				log.Fatal(err)
			}
			//err handlings**********************

			defer rows.Close() // Close the rows after the query is finished

			var curr_carttable []*models.Cart
			for rows.Next() { // Loop through the rows and add the cart of the customer to the curr_carttable slice
				cart := &models.Cart{}
				err := rows.Scan(&cart.ID, &cart.ProductID, &cart.Quantity, &cart.Price, &cart.Discount)

				//err handlings**********************
				if err != nil {
					log.Fatal(err)
				}
				//err handlings**********************

				curr_carttable = append(curr_carttable, cart)
			}

			//err handlings**********************
			if err = rows.Err(); err != nil {
				log.Fatal(err)
			}
			//err handlings**********************

			for index, element := range curr_carttable { // Loop through the curr_carttable slice and update the cart of the customer with the new price

				var product float64

				err := db.DB.QueryRow("SELECT price FROM products WHERE id = $1", element.ProductID).Scan(&product) // SQL query to get the original price of the product from the database

				//err handlings**********************
				switch {
				case err == sql.ErrNoRows:
					fmt.Println("No rows found")
					return
				case err != nil:
					log.Fatal(err)
				}
				//err handlings**********************

				curr_carttable[index].Price = totalPrice(curr_carttable[index].Quantity, product) // Update the price of the product in the cart of the customer
			}
			var total_price float64
			for _, element := range curr_carttable { // Loop through the curr_carttable to update all of the cart items of the customer
				total_price += element.Price
			}
			result, err := db.DB.Exec("INSERT INTO order_details(customer_id,total_price) VALUES($1,$2)", order_details.CustomerID, total_price) // SQL query to insert the order details of the customer to the database

			//err handlings**********************
			if err != nil {
				log.Fatal(err)
			}
			count, err := result.RowsAffected()
			if err != nil {
				log.Fatal(count)
			}
			//err handlings**********************

			table, err := db.DB.Exec("DELETE FROM cart") // Delete the cart items after the order is placed

			//err handlings**********************
			if err != nil {
				log.Fatal(err)
			}
			effect, err := table.RowsAffected()
			if err != nil {
				log.Fatal(effect)
			}
			//err handlings**********************

			w.Write([]byte("You made your order, Thanks! You just can use the Cart Quantity promotion."))

		} else if len(count_Arr) >= 1 { // If the length of the count_Arr is greater than or equal to 1, than the customer has already placed an order more than 100 so the other discounts can applied
			last_index := len(count_Arr) // Get the last index of the count_Arr
			if (last_index%4+1)%4 == 0 { // If the number of purchases is divisible by 4
				rows, err := db.DB.Query("SELECT * FROM cart") // SQL query to get the cart of the customer from the database

				//err handlings**********************
				if err != nil {
					if err == sql.ErrNoRows {
						fmt.Println("No rows found")
						return
					}
					log.Fatal(err)
				}
				//err handlings**********************

				defer rows.Close() // Close the rows after the query is finished

				var curr_carttable []*models.Cart
				for rows.Next() { // Loop through the rows and add the cart of the customer to the curr_carttable slice
					cart := &models.Cart{}
					err := rows.Scan(&cart.ID, &cart.ProductID, &cart.Quantity, &cart.Price, &cart.Discount)
					if err != nil {
						log.Fatal(err)
					}
					curr_carttable = append(curr_carttable, cart)
				}

				//err handlings**********************
				if err = rows.Err(); err != nil {
					log.Fatal(err)
				}
				//err handlings**********************

				for index, element := range curr_carttable { // Loop through the curr_carttable slice and update the cart of the customer with the new price

					var product float64
					var tax float64

					err := db.DB.QueryRow("SELECT price,tax FROM products WHERE id = $1", element.ProductID).Scan(&product, &tax) // getting tax and price of the product

					//err handlings**********************
					switch {
					case err == sql.ErrNoRows:
						fmt.Println("No rows found")
						return
					case err != nil:
						log.Fatal(err)
					}
					//err handlings**********************

					if tax == 1.00 { // If the tax is 1.00, then the product is not discountable.
						curr_carttable[index].Price = product * float64(curr_carttable[index].Quantity)
					} else if tax == 8.00 { // If the tax is 8.00, then the product is discountable by %10.
						curr_carttable[index].Price = (product - (product * 0.1)) * float64(curr_carttable[index].Quantity)
					} else if tax == 18.00 { // If the tax is 18.00, then the product is discountable by %15.
						curr_carttable[index].Price = (product - (product * 0.15)) * float64(curr_carttable[index].Quantity)
					}
				}
				var total_price float64
				for _, element := range curr_carttable { // Loop through the curr_carttable to update all of the cart item's price of the customer

					total_price += element.Price

				}
				result, err := db.DB.Exec("INSERT INTO order_details(customer_id,total_price) VALUES($1,$2)", order_details.CustomerID, total_price) // SQL query to insert the order details of the customer to the database

				//err handlings**********************
				if err != nil {
					log.Fatal(err)
				}
				count, err := result.RowsAffected()
				if err != nil {
					log.Fatal(count)
				}
				//err handlings**********************

				table, err := db.DB.Exec("DELETE FROM cart") // Delete the cart items after the order is placed

				//err handlings**********************
				if err != nil {
					log.Fatal(err)
				}
				effect, err := table.RowsAffected()
				if err != nil {
					log.Fatal(effect)
				}
				//err handlings**********************

				w.Write([]byte("You made your 4th order of higher than given amount, Thanks! The discounts of taxes are applied to all products in your cart."))

			} else { // If the number of purchases is not divisible by 4, then the customer can only use the %10 discount.
				rows, err := db.DB.Query("SELECT * FROM cart WHERE total_discount > 0") // SQL query to get the cart of the customer from the database to check that if the customer used cart Quantity promotion

				//err handlings**********************
				if err != nil {
					if err == sql.ErrNoRows {
						fmt.Println("No rows found")
						return
					}
					log.Fatal(err)
				}
				//err handlings**********************

				defer rows.Close() // Close the rows after the query is finished

				var curr_carttable []*models.Cart
				for rows.Next() { // Loop through the rows and add the cart of the customer to the curr_carttable slice
					cart := &models.Cart{}
					err := rows.Scan(&cart.ID, &cart.ProductID, &cart.Quantity, &cart.Price, &cart.Discount)
					if err != nil {
						log.Fatal(err)
					}
					curr_carttable = append(curr_carttable, cart)
				}

				//err handlings**********************
				if err = rows.Err(); err != nil {
					log.Fatal(err)
				}
				//err handlings**********************

				var discount_Arr []float64

				for _, element := range curr_carttable { // Loop through the curr_carttable slice to reach total Quantity promotion of the customer
					if element.Discount > 0 {
						discount_Arr = append(discount_Arr, element.Discount)
					}
				}
				if len(discount_Arr) < 1 { // If the customer did not use the Quantity promotion, then the customer can use the %10 discount.
					var curr_carttable []*models.Cart
					for rows.Next() { // Loop through the rows and add the cart of the customer to the curr_carttable slice
						cart := &models.Cart{}
						err := rows.Scan(&cart.ID, &cart.ProductID, &cart.Quantity, &cart.Price, &cart.Discount)
						if err != nil {
							log.Fatal(err)
						}
						curr_carttable = append(curr_carttable, cart)
					}

					//err handlings**********************
					if err = rows.Err(); err != nil {
						log.Fatal(err)
					}
					//err handlings**********************

					for index, element := range curr_carttable { // Loop through the curr_carttable slice to update the cart item's price of the customer without the discount

						var product float64

						err := db.DB.QueryRow("SELECT price FROM products WHERE id = $1", element.ProductID).Scan(&product) // Get the price of the product
						switch {
						case err == sql.ErrNoRows:
							fmt.Println("No rows found")
							return
						case err != nil:
							log.Fatal(err)
						}

						curr_carttable[index].Price = product * float64(curr_carttable[index].Quantity)
					}
					var total_price float64
					for _, element := range curr_carttable { // Loop through the curr_carttable to get total_price of the customer's cart
						total_price += element.Price
					}
					discount := total_price * 0.1                                                                                                                 // Calculate the discount of the customer
					result, err := db.DB.Exec("INSERT INTO order_details(customer_id,total_price) VALUES($1,$2)", order_details.CustomerID, total_price-discount) // SQL query to insert the order details of the customer to the database with %10 discount

					//err handlings**********************
					if err != nil {
						log.Fatal(err)
					}
					count, err := result.RowsAffected()
					if err != nil {
						log.Fatal(count)
					}
					//err handlings**********************

					table, err := db.DB.Exec("DELETE FROM cart") // Delete the cart items after the order is placed

					//err handlings**********************
					if err != nil {
						log.Fatal(err)
					}
					effect, err := table.RowsAffected()
					if err != nil {
						log.Fatal(effect)
					}
					//err handlings**********************

					w.Write([]byte("You made your order higher than given amount, Thanks! Your cart has not cart Quantity promotion. %10 discount is applied to total price of your cart."))

				} else if len(discount_Arr) >= 1 { // If the customer used the Quantity promotion, reset the cart of the customer and make the customer to use the %10 discount.

					rows, err := db.DB.Query("SELECT * FROM cart") // SQL Query to get the cart of the customer from the database to reset the cart of the customer

					//err handlings**********************
					if err != nil {
						if err == sql.ErrNoRows {
							fmt.Println("No rows found")
							return
						}
						log.Fatal(err)
					}
					//err handlings**********************

					defer rows.Close() // Close the rows after the query is finished

					var curr_carttable []*models.Cart
					for rows.Next() { // Loop through the rows and add the cart of the customer to the curr_carttable slice
						cart := &models.Cart{}
						err := rows.Scan(&cart.ID, &cart.ProductID, &cart.Quantity, &cart.Price, &cart.Discount)
						if err != nil {
							log.Fatal(err)
						}
						curr_carttable = append(curr_carttable, cart)
					}

					//err handlings**********************
					if err = rows.Err(); err != nil {
						log.Fatal(err)
					}
					//err handlings**********************

					for index, element := range curr_carttable { // Loop through the curr_carttable slice to reset the cart item's price of the customer without the discount

						var product float64

						err := db.DB.QueryRow("SELECT price FROM products WHERE id = $1", element.ProductID).Scan(&product) // Get the price of the product
						switch {
						case err == sql.ErrNoRows:
							fmt.Println("No rows found")
							return
						case err != nil:
							log.Fatal(err)
						}

						curr_carttable[index].Price = product * float64(curr_carttable[index].Quantity)
					}
					var total_price float64
					for _, element := range curr_carttable { // Loop through the curr_carttable to get total_price of the customer's cart
						total_price += element.Price
					}
					discount := total_price * 0.1                                                                                                                 // 10% discount
					result, err := db.DB.Exec("INSERT INTO order_details(customer_id,total_price) VALUES($1,$2)", order_details.CustomerID, total_price-discount) // SQL query to insert the order details of the customer to the database with %10 discount

					//err handlings**********************
					if err != nil {
						log.Fatal(err)
					}
					count, err := result.RowsAffected()
					if err != nil {
						log.Fatal(count)
					}
					//err handlings**********************

					table, err := db.DB.Exec("DELETE FROM cart") // Delete the cart items after the order is placed

					//err handlings**********************
					if err != nil {
						log.Fatal(err)
					}
					effect, err := table.RowsAffected()
					if err != nil {
						log.Fatal(effect)
					}
					//err handlings**********************

					w.Write([]byte("You made your order higher than given amount, Thanks! Your cart has cart Quantity promotion but it is declined. %10 discount is applied to total price of your cart."))

				}
			}
		}
	}
}
