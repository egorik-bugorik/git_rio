package api

import (
	"net/http"
	"pgxrio/internal/inventory"
	"strings"
)

func NewHttpServer(service *inventory.Service) http.Handler {

	s := HTTPServer{
		service: service,
		mux:     http.NewServeMux(),
	}
	s.mux.HandleFunc("/product", s.handlerGetProducts)
	return s.mux
}

type HTTPServer struct {
	service *inventory.Service
	mux     *http.ServeMux
}

func (s *HTTPServer) handlerGetProducts(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/product"):]
	if id == "" || strings.ContainsRune(id, '/') {
		http.NotFound(w, r)
		return
	}
	product,err:= s.service.
}
