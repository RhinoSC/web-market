package application

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rhinosc/web-market/code/internal"
	"github.com/rhinosc/web-market/code/internal/handler"
	"github.com/rhinosc/web-market/code/internal/repository"
	"github.com/rhinosc/web-market/code/internal/service"
)

type DefaultHTTP struct {
	addr string
}

func NewDefaultHTTP(addr string) *DefaultHTTP {
	return &DefaultHTTP{
		addr: addr,
	}
}

func (d *DefaultHTTP) Run() (err error) {
	rp := repository.NewProductRepository(make(map[int]internal.Product))

	sv := service.NewProductDefault(rp)

	hd := handler.NewDefaultProducts(sv)

	rt := chi.NewRouter()

	rt.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	rt.Route("/products", func(r chi.Router) {
		r.Get("/", hd.GetAll())
		r.Get("/{id}", hd.GetByID())
	})

	//run http server
	err = http.ListenAndServe(d.addr, rt)
	return
}
