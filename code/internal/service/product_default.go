package service

import (
	"fmt"
	"regexp"
	"time"

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

func (p *ProductDefault) GetAll() (products map[int]*internal.Product, err error) {
	return p.rp.GetAll()
}

func (p *ProductDefault) GetByID(id int) (product *internal.Product, err error) {
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

func (p *ProductDefault) SearchByPrice(price float64) (products map[int]*internal.Product, err error) {
	allProducts, err := (*p).GetAll()
	products = make(map[int]*internal.Product)
	if err != nil {
		return
	}
	for _, product := range allProducts {
		if product.Price >= price {
			products[product.Id] = product
		}
	}
	if len(products) == 0 {
		err = fmt.Errorf("%w: price", internal.ErrProductNotFound)
	}
	return
}

func (p *ProductDefault) Create(product *internal.Product) (err error) {

	if err = Validate(product); err != nil {
		return
	}

	err = p.rp.Create(product)
	return
}

func Validate(p *internal.Product) (err error) {
	// required fields
	if (*p).Name == "" {
		err = fmt.Errorf("%w: name", internal.ErrFieldRequired)
		return
	}
	if (*p).Code_value == "" {
		err = fmt.Errorf("%w: code_value", internal.ErrFieldRequired)
		return
	}
	if (*p).Expiration.IsZero() {
		err = fmt.Errorf("%w: expiration", internal.ErrFieldRequired)
		return
	}

	// quality fields
	if p.Quantity < 0 {
		err = fmt.Errorf("%w: quantity", internal.ErrValidateQualityField)
		return
	}

	rx := regexp.MustCompile(`^[A-Z][0-9]{1,5}[A-Z]?$`)
	if !rx.MatchString((*p).Code_value) {
		err = fmt.Errorf("%w: code_value", internal.ErrValidateQualityField)
		return
	}

	if p.Expiration.Before(time.Now()) {
		err = fmt.Errorf("%w: expiration", internal.ErrValidateQualityField)
		return
	}
	if p.Price < 0 {
		err = fmt.Errorf("%w: price", internal.ErrValidateQualityField)
		return
	}

	return
}
