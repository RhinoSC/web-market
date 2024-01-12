package repository

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/rhinosc/web-market/code/internal"
)

type ProductMap struct {
	db     map[int]*internal.Product
	lastID int
}

func NewProductRepository(db map[int]*internal.Product, lastID int) *ProductMap {
	pMap := &ProductMap{
		db:     db,
		lastID: lastID,
	}
	pMap.ReadProducts()
	return pMap
}

func (p *ProductMap) GetAll() (products map[int]*internal.Product, err error) {
	return p.db, nil
}

func (p *ProductMap) GetByID(id int) (product *internal.Product, err error) {
	product, ok := p.db[id]
	if !ok {
		err = fmt.Errorf("%w: id", internal.ErrProductNotFound)
		return
	}

	return
}

func (p *ProductMap) Create(product *internal.Product) (err error) {
	p.lastID++
	product.Id = p.lastID
	p.db[product.Id] = product
	return
}

type ProductJSON struct {
	Id           int     `json:"id"`
	Name         string  `json:"name"`
	Quantity     int     `json:"quantity"`
	Code_value   string  `json:"code_value"`
	Is_published bool    `json:"is_published"`
	Expiration   string  `json:"expiration"`
	Price        float64 `json:"price"`
}

func (p *ProductMap) ReadProducts() {
	// function to read products from products.json file and create a slice of products and then convert it to a map
	f, err := os.Open("products.json")
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

	var products []ProductJSON
	err = json.NewDecoder(f).Decode(&products)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, v := range products {
		t, err := time.Parse("02/01/2006", v.Expiration)
		if err != nil {
			fmt.Println(err)
		}
		(*p).db[v.Id] = &internal.Product{
			Id:           v.Id,
			Name:         v.Name,
			Quantity:     v.Quantity,
			Code_value:   v.Code_value,
			Is_published: v.Is_published,
			Expiration:   t,
			Price:        v.Price,
		}
		p.lastID = v.Id
	}
}
