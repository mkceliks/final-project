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

func GetCustomers(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query("SELECT * FROM customer")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var customers []*models.Customer
	for rows.Next() {
		var customer models.Customer
		err := rows.Scan(&customer.ID, &customer.UserName)
		if err != nil {
			log.Fatal(err)
		}
		customers = append(customers, &customer)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(customers)
}

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
		err := rows.Scan(&prd.ID, &prd.Name, &prd.Description, &prd.Price, &prd.Tax)
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
	err := row.Scan(&prd.ID, &prd.Name, &prd.Description, &prd.Price, &prd.Tax)
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

		result, err := db.DB.Exec("INSERT INTO products(name, description, price, tax) VALUES($1, $2, $3, $4)", prd.Name, prd.Description, prd.Price, prd.Tax)
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

func InsertCustomer(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		w.Header().Set("Content-Type", "application/json")
		var customer models.Customer
		_ = json.NewDecoder(r.Body).Decode(&customer)

		result, err := db.DB.Exec("INSERT INTO customer(username) VALUES($1)", customer.UserName)
		if err != nil {
			log.Fatal(err)
		}
		count, err := result.RowsAffected()
		if err != nil {
			log.Fatal(count)
		}
		w.Write([]byte("Customer added successfully"))

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

		var tax float64
		err = db.DB.QueryRow("SELECT tax FROM products WHERE id = $1", id).Scan(&tax) // Get the tax of the product
		switch {
		case err == sql.ErrNoRows:
			fmt.Println("No rows found")
			return
		case err != nil:
			log.Fatal(err)
		}

		if db.DB.QueryRow("SELECT * FROM cart WHERE product_id = $1", id).Scan(&cart.ProductID) == sql.ErrNoRows && cart.Quantity > 3 { // Check if the product is already in the cart
			floatVar, _ := strconv.ParseFloat(product, 64)
			result, err := db.DB.Exec("INSERT INTO cart(product_id, quantity, total_price, total_discount) VALUES($1, $2, $3, $4)", id, cart.Quantity, totalPrice(cart.Quantity, floatVar), totalDiscount(cart.Quantity, floatVar))
			if err != nil {
				log.Fatal(err)
			}
			count, err := result.RowsAffected()
			if err != nil {
				log.Fatal(count)
			}
			w.Write([]byte("HİÇ YOK VE 3 TEN BUYUK SAYI SEPETE EKLENDİ"))
		} else if db.DB.QueryRow("SELECT * FROM cart WHERE product_id = $1", id).Scan(&cart.ProductID) == sql.ErrNoRows && cart.Quantity <= 3 {

			floatVar, _ := strconv.ParseFloat(product, 64) // 8% discount

			result, err := db.DB.Exec("INSERT INTO cart(product_id, quantity, total_price, total_discount) VALUES($1, $2, $3, $4)", id, cart.Quantity, totalPrice(cart.Quantity, floatVar), totalDiscount(cart.Quantity, floatVar))

			if err != nil {
				log.Fatal(err)
			}
			count, err := result.RowsAffected()
			if err != nil {
				log.Fatal(count)
			}
			w.Write([]byte("HİÇ YOK VE 3 TEN AZ SAYI SEPETE EKLENDI"))
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
			floatVar, _ := strconv.ParseFloat(product, 64)
			newQuantity := cart.Quantity + oldQuantity
			result, err := db.DB.Exec("UPDATE cart SET quantity = $1, total_price = $2, total_discount = $3 WHERE product_id = $4", newQuantity, totalPrice(newQuantity, floatVar), totalDiscount(newQuantity, floatVar), id)
			if err != nil {
				log.Fatal(err)
			}

			count, err := result.RowsAffected()
			if err != nil {
				log.Fatal(count)
			}
			w.Write([]byte("SEPETTE VAR VE SEPET GUNCELLENDI"))
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

	var tax float64
	err = db.DB.QueryRow("SELECT tax FROM products WHERE id = $1", id).Scan(&tax) // Get the tax of the product
	switch {
	case err == sql.ErrNoRows:
		fmt.Println("No rows found")
		return
	case err != nil:
		log.Fatal(err)
	}

	if db.DB.QueryRow("SELECT * FROM cart WHERE product_id = $1", id).Scan(&cart.ProductID) == sql.ErrNoRows { // Check if the product is already in the cart
		floatVar, _ := strconv.ParseFloat(product, 64)
		result, err := db.DB.Exec("INSERT INTO cart(product_id, quantity, total_price, total_discount) VALUES($1, $2 + 1 , $3, $4)", id, cart.Quantity, floatVar, totalDiscount(cart.Quantity, floatVar))
		if err != nil {
			log.Fatal(err)
		}
		count, err := result.RowsAffected()
		if err != nil {
			log.Fatal(count)
		}
		w.Write([]byte("HİÇ YOK İSE 1 EKLENDI"))
	} else {
		err := db.DB.QueryRow("SELECT quantity FROM cart WHERE product_id = $1", id).Scan(&cart.Quantity)
		switch {
		case err == sql.ErrNoRows:
			fmt.Println("No rows found")
			return
		case err != nil:
			log.Fatal(err)
		}

		if db.DB.QueryRow("SELECT * FROM cart WHERE product_id = $1", id).Scan(&cart.ProductID) != sql.ErrNoRows && cart.Quantity >= 3 { // Check if the product is already in the cart
			floatVar, _ := strconv.ParseFloat(product, 64)
			newQuantity := cart.Quantity + 1 //total price of the product
			result, err := db.DB.Exec("UPDATE cart SET quantity = quantity + 1, total_price = $1, total_discount = $2 WHERE product_id = $3", totalPrice(newQuantity, floatVar), totalDiscount(newQuantity, floatVar), id)
			if err != nil {
				log.Fatal(err)
			}
			count, err := result.RowsAffected()
			if err != nil {
				log.Fatal(count)
			}
			w.Write([]byte("QUANTITY 3 VE DAHA FAZLA İSE 1 ARTTIRILDI"))
			return
		} else if db.DB.QueryRow("SELECT * FROM cart WHERE product_id = $1", id).Scan(&cart.ProductID) != sql.ErrNoRows && cart.Quantity < 3 {
			floatVar, _ := strconv.ParseFloat(product, 64)
			err := db.DB.QueryRow("SELECT quantity FROM cart WHERE product_id = $1", id).Scan(&cart.Quantity)
			switch {
			case err == sql.ErrNoRows:
				fmt.Println("No rows found")
				return
			case err != nil:
				log.Fatal(err)
			}
			newQuantity := cart.Quantity + 1
			result, err := db.DB.Exec("UPDATE cart SET quantity = quantity + 1, total_price = $1, total_discount = $2 WHERE product_id = $3", totalPrice(newQuantity, floatVar), totalDiscount(newQuantity, floatVar), id)
			if err != nil {
				log.Fatal(err)
			}
			count, err := result.RowsAffected()
			if err != nil {
				log.Fatal(count)
			}
			w.Write([]byte("QUANTITY 3den KÜÇÜKSE ve var ise 1 ARTTIRILDI"))
		}
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
		} else if newQuantity > 3 {
			//total price of the product
			result, err := db.DB.Exec("UPDATE cart SET quantity = $1, total_price = $2, total_discount = $3 WHERE product_id = $4", newQuantity, totalPrice(newQuantity, floatVar), totalDiscount(newQuantity, floatVar), id)
			if err != nil {
				log.Fatal(err)
			}
			count, err := result.RowsAffected()
			if err != nil {
				log.Fatal(count)
			}
		} else {
			result, err := db.DB.Exec("UPDATE cart SET quantity = $1, total_price = $2, total_discount = $3 WHERE product_id = $4", newQuantity, totalPrice(newQuantity, floatVar), totalDiscount(newQuantity, floatVar), id)
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

func Order(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		w.Header().Set("Content-Type", "application/json")
		var order_details models.OrderDetails
		_ = json.NewDecoder(r.Body).Decode(&order_details)

		rows, err := db.DB.Query("SELECT * FROM order_details WHERE customer_id = $1", order_details.CustomerID)
		if err != nil {
			if err == sql.ErrNoRows {
				fmt.Println("No rows found")
				return
			}
			log.Fatal(err)
		}
		defer rows.Close()

		var curr_ordertable []*models.OrderDetails
		for rows.Next() {
			ord := &models.OrderDetails{}
			err := rows.Scan(&ord.ID, &ord.CustomerID, &ord.TotalPrice)
			if err != nil {
				log.Fatal(err)
			}
			curr_ordertable = append(curr_ordertable, ord)
		}
		if err = rows.Err(); err != nil {
			log.Fatal(err)
		}
		var count_Arr []float32
		for _, element := range curr_ordertable {
			if element.TotalPrice > 100 {
				count_Arr = append(count_Arr, element.TotalPrice)
			}
		}
		if len(count_Arr) < 1 {
			result, err := db.DB.Exec("INSERT INTO order_details(customer_id,total_price) VALUES($1,$2)", order_details.CustomerID, order_details.TotalPrice)
			if err != nil {
				log.Fatal(err)
			}
			count, err := result.RowsAffected()
			if err != nil {
				log.Fatal(count)
			}

			table, err := db.DB.Exec("DELETE FROM cart") // Delete the cart items after the order is placed
			if err != nil {
				log.Fatal(err)
			}
			effect, err := table.RowsAffected()
			if err != nil {
				log.Fatal(effect)
			}

			w.Write([]byte("Purchase successfull"))

		} else if len(count_Arr) >= 1 {

			rows, err := db.DB.Query("SELECT * FROM cart WHERE total_discount > 0")
			if err != nil {
				if err == sql.ErrNoRows {
					fmt.Println("No rows found")
					return
				}
				log.Fatal(err)
			}
			defer rows.Close()

			var curr_carttable []*models.Cart
			for rows.Next() {
				cart := &models.Cart{}
				err := rows.Scan(&cart.ID, &cart.ProductID, &cart.Quantity, &cart.Price, &cart.Discount)
				if err != nil {
					log.Fatal(err)
				}
				curr_carttable = append(curr_carttable, cart)
			}
			if err = rows.Err(); err != nil {
				log.Fatal(err)
			}
			var discount_Arr []float64

			for _, element := range curr_carttable {
				if element.Discount > 0 {
					discount_Arr = append(discount_Arr, element.Discount)
				}
			}
			if len(discount_Arr) < 1 {
				discount := order_details.TotalPrice * 0.1 // 10% discount

				result, err := db.DB.Exec("INSERT INTO order_details(customer_id,total_price) VALUES($1,$2)", order_details.CustomerID, order_details.TotalPrice-discount)
				if err != nil {
					log.Fatal(err)
				}
				count, err := result.RowsAffected()
				if err != nil {
					log.Fatal(count)
				}

				table, err := db.DB.Exec("DELETE FROM cart") // Delete the cart items after the order is placed
				if err != nil {
					log.Fatal(err)
				}
				effect, err := table.RowsAffected()
				if err != nil {
					log.Fatal(effect)
				}

				w.Write([]byte("Purchase successfull"))

			} else if len(discount_Arr) >= 1 {

				rows, err := db.DB.Query("SELECT * FROM cart")
				if err != nil {
					if err == sql.ErrNoRows {
						fmt.Println("No rows found")
						return
					}
					log.Fatal(err)
				}
				defer rows.Close()

				var curr_carttable []*models.Cart
				for rows.Next() {
					cart := &models.Cart{}
					err := rows.Scan(&cart.ID, &cart.ProductID, &cart.Quantity, &cart.Price, &cart.Discount)
					if err != nil {
						log.Fatal(err)
					}
					curr_carttable = append(curr_carttable, cart)
				}
				if err = rows.Err(); err != nil {
					log.Fatal(err)
				}

				for index, element := range curr_carttable {

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
					fmt.Printf("%+v\n", curr_carttable[index].Price)
					fmt.Printf("%d\n", index)
				}

				for _, element := range curr_carttable {
					var total_price float64
					total_price += element.Price
					fmt.Printf("%f\n", total_price)

					discount := total_price * 0.1 // 10% discount

					result, err := db.DB.Exec("INSERT INTO order_details(customer_id,total_price) VALUES($1,$2)", order_details.CustomerID, total_price-discount)
					if err != nil {
						log.Fatal(err)
					}
					count, err := result.RowsAffected()
					if err != nil {
						log.Fatal(count)
					}
				}

				table, err := db.DB.Exec("DELETE FROM cart") // Delete the cart items after the order is placed
				if err != nil {
					log.Fatal(err)
				}
				effect, err := table.RowsAffected()
				if err != nil {
					log.Fatal(effect)
				}

				w.Write([]byte("Purchase successfull"))

			}
		}
	}
}
