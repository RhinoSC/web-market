package service

import (
	"fmt"

	"github.com/rhinosc/web-market/code/internal"
)

type ProductDefault struct {
	rp internal.ProductRepository
}

func NewProductDefault(rp internal.ProductRepository) *ProductDefault {
	return &ProductDefault{
		rp: rp,
	}
}

func (p *ProductDefault) GetAll() (products map[int]internal.Product, err error) {
	return p.rp.GetAll()
}

func (p *ProductDefault) GetByID(id int) (product internal.Product, err error) {
	product, err = p.rp.GetByID(id)
	if err != nil {
		switch err {
		case internal.ErrProductNotFound:
			err = fmt.Errorf("%w: id", internal.ErrProductNotFound)
		}
		return
	}

	return
}
