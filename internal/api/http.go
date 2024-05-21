package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"pgxrio/internal/inventory"
	"strings"
)

func NewHttpServer(service *inventory.Service) http.Handler {

	s := HTTPServer{
		service: service,
		mux:     http.NewServeMux(),
	}
	s.mux.HandleFunc("/product/", s.handlerGetProducts)
	return s.mux
}

type HTTPServer struct {
	service *inventory.Service
	mux     *http.ServeMux
}

func (s *HTTPServer) handlerGetProducts(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/product/"):]
	if id == "" || strings.ContainsRune(id, '/') {
		http.NotFound(w, r)
		return
	}
	product, err := s.service.GetProduct(r.Context(), id)
	if err == context.Canceled || err == context.DeadlineExceeded {
		return
	}
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Println(err)
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "\t")
	err = enc.Encode(product)
	if err != nil {
		log.Printf("Can't json response ::: %v", err)
	}

}
