package application

import (
	"net/http"

	"github.com/go-chi/chi/v5"
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

	st := repository.NewStorageProductJSON("products1.json", "02/01/2006")
	// rp := repository.NewProductRepository(make(map[int]*internal.Product), 0)
	rp := repository.NewProductStore(*st, 0, "02/01/2006")

	sv := service.NewProductDefault(rp)

	hd := handler.NewDefaultProducts(sv)

	rt := chi.NewRouter()

	rt.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	rt.Route("/products", func(r chi.Router) {
		r.Get("/", hd.GetAll())
		r.Get("/{id}", hd.GetByID())
		r.Get("/search", hd.Search())

		r.Post("/", hd.Create())

		r.Put("/{id}", hd.UpdateOrCreate())
		r.Patch("/{id}", hd.Update())
		r.Delete("/{id}", hd.Delete())
	})

	//run http server
	err = http.ListenAndServe(d.addr, rt)
	return
}
