package models

type Cart struct {
	ID        int
	ProductID int
	Quantity  int
	Price     float64
	Discount  float64
}
