package internal

import "errors"

var (
	ErrProductNotFound = errors.New("product not found")
)

type ProductRepository interface {
	// Returns all products
	GetAll() (products map[int]*Product, err error)

	// Returns a product by ID
	GetByID(id int) (product *Product, err error)

	// Creates a new product
	Create(product *Product) (err error)
}
