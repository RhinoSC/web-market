package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/rhinosc/web-market/code/internal"
	"github.com/rhinosc/web-market/code/platform/web/response"
)

type DefaultProducts struct {
	sv internal.ProductService
}

func NewDefaultProducts(sv internal.ProductService) *DefaultProducts {
	return &DefaultProducts{
		sv: sv,
	}
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

type BodyProductJSON struct {
	Name         string  `json:"name"`
	Quantity     int     `json:"quantity"`
	Code_value   string  `json:"code_value"`
	Is_published bool    `json:"is_published"`
	Expiration   string  `json:"expiration"`
	Price        float64 `json:"price"`
}

// GetAll returns all products
func (p *DefaultProducts) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//request

		//process
		products, err := p.sv.GetAll()
		if err != nil {
			switch err {
			case internal.ErrProductNotFound:
				response.Text(w, http.StatusNotFound, "Product not found")
			default:
				response.Text(w, http.StatusInternalServerError, "Internal Server Error")
			}
			return
		}

		//response
		// serialize products to json
		var data []ProductJSON
		for _, products := range products {
			pJSON := ProductJSON{
				Id:           products.Id,
				Name:         products.Name,
				Quantity:     products.Quantity,
				Code_value:   products.Code_value,
				Is_published: products.Is_published,
				Expiration:   products.Expiration.Format("2006-01-02"),
				Price:        products.Price,
			}
			data = append(data, pJSON)
		}
		response.JSON(w, http.StatusOK, map[string]any{
			"message": "success",
			"data":    data,
		})
	}
}

// GetByID returns a product by id
func (p *DefaultProducts) GetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//request

		//get id from urlparams with chi
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			response.Text(w, http.StatusBadRequest, "invalid id")
			return
		}

		//process
		product, err := p.sv.GetByID(id)
		if err != nil {
			switch err {
			case internal.ErrProductNotFound:
				response.Text(w, http.StatusNotFound, "Product not found")
			default:
				response.Text(w, http.StatusInternalServerError, "Internal Server Error")
			}
			return
		}

		//response
		// serialize product to json
		data := ProductJSON{
			Id:           product.Id,
			Name:         product.Name,
			Quantity:     product.Quantity,
			Code_value:   product.Code_value,
			Is_published: product.Is_published,
			Expiration:   product.Expiration.Format("2006-01-02"),
			Price:        product.Price,
		}

		response.JSON(w, http.StatusOK, map[string]any{
			"message": "success",
			"data":    data,
		})
	}
}

// Search returns a filtered map of products where are greater than the given price
func (p *DefaultProducts) Search() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//request

		//get price from urlparams with chi
		price, err := strconv.ParseInt(r.URL.Query().Get("priceGt"), 10, 64)
		if err != nil {
			response.Text(w, http.StatusBadRequest, "invalid price")
			return
		}

		//process
		products, err := p.sv.SearchByPrice(float64(price))
		if err != nil {
			switch err {
			case internal.ErrProductNotFound:
				response.Text(w, http.StatusNotFound, "Product not found")
			default:
				response.Text(w, http.StatusInternalServerError, "Internal Server Error")
			}
			return
		}

		//response
		// serialize products to json
		var data []ProductJSON
		for _, products := range products {
			pJSON := ProductJSON{
				Id:           products.Id,
				Name:         products.Name,
				Quantity:     products.Quantity,
				Code_value:   products.Code_value,
				Is_published: products.Is_published,
				Expiration:   products.Expiration.Format("2006-01-02"),
				Price:        products.Price,
			}
			data = append(data, pJSON)
		}
		response.JSON(w, http.StatusOK, map[string]any{
			"message": "success",
			"data":    data,
		})
	}
}

// Create creates a product
func (p *DefaultProducts) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//request

		//decode body to json
		var body BodyProductJSON
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			response.Text(w, http.StatusBadRequest, "invalid body")
			return
		}

		//process
		exp, err := time.Parse("02/01/2006", body.Expiration)
		if err != nil {
			code := http.StatusBadRequest
			response.JSON(w, code, map[string]any{
				"message": "invalid expiration",
				"data":    nil,
			})
			return
		}

		product := internal.Product{
			Name:         body.Name,
			Quantity:     body.Quantity,
			Code_value:   body.Code_value,
			Is_published: body.Is_published,
			Expiration:   exp,
			Price:        body.Price,
		}
		err = p.sv.Create(&product)
		if err != nil {
			switch {
			case errors.Is(err, internal.ErrFieldRequired):
				response.Text(w, http.StatusBadRequest, "Field required")
			case errors.Is(err, internal.ErrValidateQualityField):
				response.Text(w, http.StatusBadRequest, "Invalid expiration")
			default:
				response.Text(w, http.StatusInternalServerError, "Internal Server Error")
			}
			return
		}

		//response
		response.JSON(w, http.StatusCreated, map[string]any{
			"message": "success",
			"data":    product,
		})
	}
}
