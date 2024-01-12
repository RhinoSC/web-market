package repository

import (
	"fmt"

	"github.com/rhinosc/web-market/code/internal"
)

type ProductStore struct {
	st StorageProductJSON

	LastID int

	LayoutDate string
}

func NewProductStore(st StorageProductJSON, lastID int, layoutDate string) *ProductStore {
	return &ProductStore{
		st:         st,
		LastID:     lastID,
		LayoutDate: layoutDate,
	}
}

func (p *ProductStore) GetAll() (products map[int]*internal.Product, err error) {
	prod, err := p.st.ReadAll()
	if err != nil {
		err = fmt.Errorf("%w: id", internal.ErrProductNotFound)
	}
	products = prod
	return
}

func (p *ProductStore) GetByID(id int) (product *internal.Product, err error) {
	prods, err := p.st.ReadAll()
	if err != nil {
		return
	}
	product, ok := prods[id]
	if !ok {
		err = fmt.Errorf("%w: id", internal.ErrProductNotFound)
		return
	}
	return
}

func (p *ProductStore) SearchByPrice(price float64) (products map[int]*internal.Product, err error) {
	prods, err := p.st.ReadAll()
	if err != nil {
		return
	}

	filteredProducts := make(map[int]*internal.Product)
	for _, v := range prods {
		if v.Price >= price {
			filteredProducts[v.Id] = v
		}
	}
	products = filteredProducts
	return
}

func (p *ProductStore) Create(product *internal.Product) (err error) {
	prods, err := p.st.ReadAll()
	if err != nil {
		return
	}
	p.LastID++
	product.Id = p.LastID
	prods[product.Id] = product

	err = p.st.WriteAll(prods)
	if err != nil {
		return
	}
	return
}

func (p *ProductStore) UpdateOrCreate(product *internal.Product) (prod internal.Product, err error) {
	prods, err := p.st.ReadAll()
	if err != nil {
		return
	}
	_, ok := prods[product.Id]
	switch ok {
	case true:
		//update
		prods[product.Id] = product
	case false:
		//create
		p.LastID++
		product.Id = p.LastID
		prods[product.Id] = product
	}
	prod = *product
	err = p.st.WriteAll(prods)
	if err != nil {
		return
	}
	return
}

func (p *ProductStore) Update(product *internal.Product) (err error) {
	prods, err := p.st.ReadAll()
	if err != nil {
		return
	}
	_, ok := prods[product.Id]
	if !ok {
		err = fmt.Errorf("%w: id", internal.ErrProductNotFound)
		return
	}
	prods[product.Id] = product
	err = p.st.WriteAll(prods)
	if err != nil {
		return
	}
	return
}

func (p *ProductStore) Delete(id int) (err error) {
	prods, err := p.st.ReadAll()
	if err != nil {
		return
	}
	_, ok := prods[id]
	if !ok {
		err = fmt.Errorf("%w: id", internal.ErrProductNotFound)
		return
	}
	delete(prods, id)
	err = p.st.WriteAll(prods)
	if err != nil {
		return
	}
	return
}
