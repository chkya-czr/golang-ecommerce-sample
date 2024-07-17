package handler

import (
	"product_service/internal/domain/product/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

func RegisterHTTPEndPoints(router *chi.Mux, validator *validator.Validate, svc service.IProductService) *Handler {
	h := NewHandler(svc, validator)

	router.Route("/api/v1/product", func(router chi.Router) {
		router.Get("/static/{productID}", h.GetStatic)
	})
	return h
}
