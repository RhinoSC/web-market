package internal

import "errors"

var (
	ErrFieldRequired        = errors.New("field required")
	ErrValidateQualityField = errors.New("validate quality field")
)

type ProductService interface {
	// Returns all products
	GetAll() (products map[int]*Product, err error)

	// Returns a product by ID
	GetByID(id int) (product *Product, err error)

	// Returns a product by price
	SearchByPrice(price float64) (products map[int]*Product, err error)

	// Creates a new product
	Create(product *Product) (err error)

	// Updates or creates a product
	UpdateOrCreate(product *Product) (prod Product, err error)

	// Updates a product
	Update(product *Product) (err error)

	// Deletes a product
	Delete(id int) (err error)
}
