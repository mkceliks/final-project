package service

var given_amount float64 // given amount of the product.

func totalPrice(quantity int, productPrice float64) float64 { // total_price = (product*3) + (product*(quantity-3)) - (discount * (quantity-3))

	if quantity > 3 { // if the quantity is greater than 3 then discount is applied.
		total_price := (productPrice * 3) + (productPrice*(float64(quantity)-3) - totalDiscount(quantity, productPrice)) // calculating the total price of the product with subtraction of discount.
		return total_price
	} else {
		total_price := (productPrice * float64(quantity)) // calculating the total price of the product.
		return total_price
	}
}

func totalDiscount(quantity int, productPrice float64) float64 { // total discount of the product.
	discount := productPrice * 0.08 // 8% discount

	if quantity > 3 { // if the quantity is greater than 3 then discount is applied.
		total_discount := (discount * (float64(quantity) - 3))
		return total_discount
	} else {
		return 0.00
	}

}

// func totalDiscountTen(quantity int, productPrice float64) float64 { // total discount of the product.
// 	discount := productPrice * 0.1 // 8% discount

// 	if quantity > 3 { // if the quantity is greater than 3 then discount is applied.
// 		total_discount := (discount * (float64(quantity) - 3))
// 		return total_discount
// 	} else {
// 		return 0.00
// 	}

// }
