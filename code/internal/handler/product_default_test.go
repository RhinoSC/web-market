package handler_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/rhinosc/web-market/code/internal"
	"github.com/rhinosc/web-market/code/internal/handler"
	"github.com/rhinosc/web-market/code/internal/repository"
	"github.com/rhinosc/web-market/code/internal/service"
	"github.com/stretchr/testify/require"
)

func TestProductDefault_GetAll(t *testing.T) {
	t.Run("success 01 - should return a list of products", func(t *testing.T) {
		// arrange
		db := make(map[int]*internal.Product)
		db[1] = &internal.Product{
			Id:           1,
			Name:         "Product 1",
			Quantity:     10,
			Code_value:   "123456",
			Is_published: true,
			Expiration:   time.Date(2006, time.February, 1, 0, 0, 0, 0, time.UTC),
			Price:        10.0,
		}

		db[2] = &internal.Product{
			Id:           2,
			Name:         "Product 2",
			Quantity:     20,
			Code_value:   "123456",
			Is_published: true,
			Expiration:   time.Date(2006, time.February, 1, 0, 0, 0, 0, time.UTC),
			Price:        20.0,
		}

		rp := repository.NewProductRepository(db, 0)
		sv := service.NewProductDefault(rp)
		hd := handler.NewDefaultProducts(sv)

		hdFunc := hd.GetAll()

		// act

		req := httptest.NewRequest("GET", "/products", nil)
		res := httptest.NewRecorder()

		hdFunc(res, req)

		type Response struct {
			Data    []handler.ProductJSON `json:"data"`
			Message string                `json:"message"`
		}

		var response Response
		json.NewDecoder(res.Body).Decode(&response)

		body, err := json.Marshal(response.Data)
		if err != nil {
			t.Fatalf("Error converting 'data' to JSON: %v", err)
		}

		// assert

		expectedCode := http.StatusOK
		expectedBody := `[{"id":1,"name":"Product 1","quantity":10,"code_value":"123456","is_published":true,"expiration":"01/02/2006","price":10},{"id":2,"name":"Product 2","quantity":20,"code_value":"123456","is_published":true,"expiration":"01/02/2006","price":20}]`
		expectedHeader := "application/json"

		require.Equal(t, expectedCode, res.Code)
		require.JSONEq(t, expectedBody, string(body))
		require.Equal(t, expectedHeader, res.Header().Get("Content-Type"))
	})

	t.Run("success 02 - should return an empty list of products", func(t *testing.T) {
		// arrange
		db := make(map[int]*internal.Product)
		rp := repository.NewProductRepository(db, 0)
		sv := service.NewProductDefault(rp)
		hd := handler.NewDefaultProducts(sv)

		hdFunc := hd.GetAll()

		// act

		req := httptest.NewRequest("GET", "/products", nil)
		res := httptest.NewRecorder()

		hdFunc(res, req)

		type Response struct {
			Data    []handler.ProductJSON `json:"data"`
			Message string                `json:"message"`
		}

		var response Response
		json.NewDecoder(res.Body).Decode(&response)

		// assert

		expectedCode := http.StatusOK
		var expectedBody []handler.ProductJSON
		expectedHeader := "application/json"

		require.Equal(t, expectedCode, res.Code)
		require.Equal(t, expectedBody, response.Data)
		require.Equal(t, expectedHeader, res.Header().Get("Content-Type"))
	})

	t.Run("success 03 - should return a product by id", func(t *testing.T) {
		// arrange
		db := make(map[int]*internal.Product)
		db[1] = &internal.Product{
			Id:           1,
			Name:         "Product 1",
			Quantity:     10,
			Code_value:   "123456",
			Is_published: true,
			Expiration:   time.Date(2006, time.February, 1, 0, 0, 0, 0, time.UTC),
			Price:        10.0,
		}

		rp := repository.NewProductRepository(db, 0)
		sv := service.NewProductDefault(rp)
		hd := handler.NewDefaultProducts(sv)

		hdFunc := hd.GetByID()

		// act

		req := httptest.NewRequest("GET", "/products/1", nil)
		chiCtx := chi.NewRouteContext()
		chiCtx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
		res := httptest.NewRecorder()

		hdFunc(res, req)

		type Response struct {
			Data    handler.ProductJSON `json:"data"`
			Message string              `json:"message"`
		}

		var response Response
		json.NewDecoder(res.Body).Decode(&response)

		body, err := json.Marshal(response.Data)
		if err != nil {
			t.Fatalf("Error converting 'data' to JSON: %v", err)
		}

		// assert

		expectedCode := http.StatusOK
		expectedBody := `{"id":1,"name":"Product 1","quantity":10,"code_value":"123456","is_published":true,"expiration":"01/02/2006","price":10}`
		expectedHeader := "application/json"

		require.Equal(t, expectedCode, res.Code)
		require.JSONEq(t, expectedBody, string(body))
		require.Equal(t, expectedHeader, res.Header().Get("Content-Type"))
	})

	t.Run("success 04 - should create a product", func(t *testing.T) {
		// arrange
		db := make(map[int]*internal.Product)
		rp := repository.NewProductRepository(db, 0)
		sv := service.NewProductDefault(rp)
		hd := handler.NewDefaultProducts(sv)

		hdFunc := hd.Create()

		product := `{"name":"Product 1","quantity":10,"code_value":"S6611","is_published":true,"expiration":"01/12/2024","price":10}`

		// act

		req := httptest.NewRequest("POST", "/products", strings.NewReader(product))
		req.Header.Set("Authorization", "12345")
		res := httptest.NewRecorder()

		hdFunc(res, req)

		type Response struct {
			Data    handler.ProductJSON `json:"data"`
			Message string              `json:"message"`
		}

		var response Response
		json.NewDecoder(res.Body).Decode(&response)

		body, err := json.Marshal(response.Data)
		if err != nil {
			t.Fatalf("Error converting 'data' to JSON: %v", err)
		}

		// assert

		expectedCode := http.StatusCreated
		expectedBody := `{"id":1,"name":"Product 1","quantity":10,"code_value":"S6611","is_published":true,"expiration":"01/12/2024","price":10}`
		expectedHeader := "application/json"

		require.Equal(t, expectedCode, res.Code)
		require.JSONEq(t, expectedBody, string(body))
		require.Equal(t, expectedHeader, res.Header().Get("Content-Type"))
	})

	t.Run("success 05 - should delete a product", func(t *testing.T) {
		// arrange
		db := make(map[int]*internal.Product)
		db[1] = &internal.Product{
			Id:           1,
			Name:         "Product 1",
			Quantity:     10,
			Code_value:   "123456",
			Is_published: true,
			Expiration:   time.Date(2006, time.February, 1, 0, 0, 0, 0, time.UTC),
			Price:        10.0,
		}

		rp := repository.NewProductRepository(db, 0)
		sv := service.NewProductDefault(rp)
		hd := handler.NewDefaultProducts(sv)

		hdFunc := hd.Delete()

		// act

		req := httptest.NewRequest("DELETE", "/products/1", nil)
		req.Header.Set("Authorization", "12345")
		chiCtx := chi.NewRouteContext()
		chiCtx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
		res := httptest.NewRecorder()

		hdFunc(res, req)

		// assert

		expectedCode := http.StatusNoContent
		expectedHeader := "application/json"

		require.Equal(t, expectedCode, res.Code)
		require.Equal(t, expectedHeader, res.Header().Get("Content-Type"))
	})

	t.Run("fail 01 - should return not found when trying to get a product by id", func(t *testing.T) {
		// arrange
		db := make(map[int]*internal.Product)
		rp := repository.NewProductRepository(db, 0)
		sv := service.NewProductDefault(rp)
		hd := handler.NewDefaultProducts(sv)

		hdFunc := hd.GetByID()

		// act

		req := httptest.NewRequest("GET", "/products/1", nil)
		chiCtx := chi.NewRouteContext()
		chiCtx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
		res := httptest.NewRecorder()

		hdFunc(res, req)

		var response string
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			t.Fatalf("Error reading response body: %v", err)
		}
		response = string(bodyBytes)

		// body, err := json.Marshal(response)
		// if err != nil {
		// 	t.Fatalf("Error converting 'data' to JSON: %v", err)
		// }

		// assert

		expectedCode := http.StatusNotFound
		expectedHeader := "text/plain; charset=utf-8"

		require.Equal(t, expectedCode, res.Code)
		require.Equal(t, expectedHeader, res.Header().Get("Content-Type"))
		require.Equal(t, `Product not found`, response)
	})

	t.Run("fail 02 - should return not found when trying to delete a product", func(t *testing.T) {
		// arrange
		db := make(map[int]*internal.Product)
		rp := repository.NewProductRepository(db, 0)
		sv := service.NewProductDefault(rp)
		hd := handler.NewDefaultProducts(sv)

		hdFunc := hd.Delete()

		// act

		req := httptest.NewRequest("DELETE", "/products/1", nil)
		req.Header.Set("Authorization", "12345")
		chiCtx := chi.NewRouteContext()
		chiCtx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
		res := httptest.NewRecorder()

		hdFunc(res, req)

		var response string
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			t.Fatalf("Error reading response body: %v", err)
		}
		response = string(bodyBytes)

		// assert

		expectedCode := http.StatusNotFound
		expectedHeader := "text/plain; charset=utf-8"

		require.Equal(t, expectedCode, res.Code)
		require.Equal(t, expectedHeader, res.Header().Get("Content-Type"))
		require.Equal(t, `Product not found`, response)
	})

	t.Run("fail 03 - should return bad request when trying to get a product with invalid id", func(t *testing.T) {
		// arrange
		db := make(map[int]*internal.Product)
		rp := repository.NewProductRepository(db, 0)
		sv := service.NewProductDefault(rp)
		hd := handler.NewDefaultProducts(sv)

		hdFunc := hd.GetByID()

		// act

		req := httptest.NewRequest("GET", "/products/A1", nil)
		chiCtx := chi.NewRouteContext()
		chiCtx.URLParams.Add("id", "A1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
		res := httptest.NewRecorder()

		hdFunc(res, req)

		var response string
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			t.Fatalf("Error reading response body: %v", err)
		}
		response = string(bodyBytes)

		// assert

		expectedCode := http.StatusBadRequest
		expectedHeader := "text/plain; charset=utf-8"

		require.Equal(t, expectedCode, res.Code)
		require.Equal(t, expectedHeader, res.Header().Get("Content-Type"))
		require.Equal(t, `invalid id`, response)
	})
}
