package models

type Product struct {
	ID                int
	Name, Description string
	Price             float32
	Discount          float32
	Tax               float32
}
