package internal

type ProductService interface {
	// Returns all products
	GetAll() (products map[int]Product, err error)

	// Returns a product by ID
	GetByID(id int) (product Product, err error)
}
