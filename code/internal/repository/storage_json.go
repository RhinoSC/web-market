package repository

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/rhinosc/web-market/code/internal"
)

type StorageProductJSON struct {
	FilePath   string
	LayoutDate string
}

func NewStorageProductJSON(filePath string, layoutDate string) *StorageProductJSON {
	return &StorageProductJSON{
		FilePath:   filePath,
		LayoutDate: layoutDate,
	}
}

type ProductAttributesJSON struct {
	Id           int     `json:"id"`
	Name         string  `json:"name"`
	Quantity     int     `json:"quantity"`
	Code_value   string  `json:"code_value"`
	Is_published bool    `json:"is_published"`
	Expiration   string  `json:"expiration"`
	Price        float64 `json:"price"`
}

func (s *StorageProductJSON) ReadAll() (p map[int]*internal.Product, err error) {
	// function to read products from products.json file and create a slice of products and then convert it to a map
	f, err := os.Open(s.FilePath)
	if err != nil {
		fmt.Println("error opening file: ", s.FilePath)
	}
	defer f.Close()

	var products []ProductJSON
	err = json.NewDecoder(f).Decode(&products)
	if err != nil {
		fmt.Println("error decoding file: ", s.FilePath)
		return
	}

	p = make(map[int]*internal.Product)
	for _, v := range products {
		t, err := time.Parse(s.LayoutDate, v.Expiration)
		if err != nil {
			fmt.Println("error parsing time: ", v.Expiration)
		}
		p[v.Id] = &internal.Product{
			Id:           v.Id,
			Name:         v.Name,
			Quantity:     v.Quantity,
			Code_value:   v.Code_value,
			Is_published: v.Is_published,
			Expiration:   t,
			Price:        v.Price,
		}
	}
	return
}

func (s *StorageProductJSON) WriteAll(p map[int]*internal.Product) (err error) {
	// function to write products to products.json file
	f, err := os.OpenFile(s.FilePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	var products []ProductJSON
	for _, v := range p {
		products = append(products, ProductJSON{
			Id:           v.Id,
			Name:         v.Name,
			Quantity:     v.Quantity,
			Code_value:   v.Code_value,
			Is_published: v.Is_published,
			Expiration:   v.Expiration.Format(s.LayoutDate),
			Price:        v.Price,
		})
	}

	err = json.NewEncoder(f).Encode(products)
	if err != nil {
		fmt.Println(err)
		return
	}
	return
}
