package handler

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"product_service/internal/domain/product/model"
	"product_service/internal/domain/product/service"
	"product_service/utility/message"
	"product_service/utility/param"
	"product_service/utility/respond"

	"github.com/go-playground/validator/v10"
)

type Handler struct {
	service  service.IProductService
	validate *validator.Validate
}

func NewHandler(service service.IProductService, validate *validator.Validate) *Handler {
	return &Handler{
		service:  service,
		validate: validate,
	}
}

// Get a product by its ID
// @Summary Get Product Info
// @Description Get a product by its id.
// @Accept json
// @Produce json
// @Param productID path int true "book ID"
// @Success 200 {object} product.Res
// @Failure 400 {string} Bad product.CreateRequest
// @Failure 500 {string} Internal Server Error
// @router /api/v1/product/{productID} [get]
func (h *Handler) GetStatic(w http.ResponseWriter, r *http.Request) {
	productId, err := param.Int(r, "productID")
	if err != nil {
		respond.Error(w, http.StatusBadRequest, message.ErrBadRequest)
		return
	}

	b, err := h.service.Get(context.Background(), productId)
	if err != nil {
		if err == sql.ErrNoRows {
			respond.Error(w, http.StatusBadRequest, errors.New("no book is found for this ID"))
			return
		}
		respond.Error(w, http.StatusInternalServerError, nil)
		return
	}
	product := model.Resource(b)

	respond.Json(w, http.StatusOK, product)
}
