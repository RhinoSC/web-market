package repository

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/rhinosc/web-market/code/internal"
)

type ProductMap struct {
	db map[int]internal.Product
}

func NewProductRepository(db map[int]internal.Product) *ProductMap {
	pMap := &ProductMap{
		db: db,
	}
	pMap.ReadProducts()
	return pMap
}

func (p *ProductMap) GetAll() (products map[int]internal.Product, err error) {
	return p.db, nil
}

func (p *ProductMap) GetByID(id int) (product internal.Product, err error) {
	product, ok := p.db[id]
	if !ok {
		err = fmt.Errorf("%w: id", internal.ErrProductNotFound)
		return
	}

	return
}

func (p *ProductMap) ReadProducts() {
	// function to read products from products.json file and create a slice of products and then convert it to a map
	f, err := os.Open("products.json")
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

	var products []internal.Product
	json.NewDecoder(f).Decode(&products)

	for _, product := range products {
		p.db[product.Id] = product
	}
}
