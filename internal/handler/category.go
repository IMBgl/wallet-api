package handler

import (
	"context"
	"errors"
	"github.com/IMBgl/go-wallet-api/internal/domain"
	"github.com/IMBgl/go-wallet-api/internal/service"
	"github.com/IMBgl/go-wallet-api/pkg/validator"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"log"
	"net/http"
)

type CategoryHandler struct {
	categoryService service.CategoryService
	middleware      *apiMiddleware
}

func (h CategoryHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Use(h.middleware.Auth)
	r.Post("/", h.create)

	return r
}

type CategoryCreateRequest struct {
	Name        string          `json:"name"`
	Currency    string          `json:"currency"`
	ParentId    *string         `json:"parentId,omitempty"`
	ParentIdVal *uuid.UUID      `json:"-"`
	CurrencyVal domain.Currency `json:"-"`
}

type CategoryResponse struct {
	Id        string              `json:"id"`
	Name      string              `json:"name"`
	Currency  string              `json:"currency"`
	UserId    string              `json:"userId"`
	CreatedAt string              `json:"createdAt"`
	ParentId  *uuid.UUID          `json:"parentId"`
	Parent    *CategoryResponse   `json:"parent"`
	Children  []*CategoryResponse `json:"children"`
}

func NewCategoryResponse(c *domain.Category, p *domain.Category) *CategoryResponse {
	resp := &CategoryResponse{
		Id:        c.Id.String(),
		UserId:    c.UserId.String(),
		Currency:  c.Currency.Val(),
		Name:      c.Name,
		CreatedAt: c.CreatedAt.Format(DateTimeFormat()),
		ParentId:  c.ParentId,
	}

	if p != nil {
		resp.Parent = NewCategoryResponse(p, nil)
	}

	return resp
}

func (data *CategoryCreateRequest) Bind(r *http.Request) error {
	if data.Name == "" {
		return errors.New("name field required")
	}

	currency, err := domain.CurrencyFromString(data.Currency)
	if err != nil {
		return errors.New("currency format must be one of 'rur', 'eur, 'usd'")
	}
	data.CurrencyVal = currency

	if data.ParentId != nil {
		parentId, err := validator.Uuid(*data.ParentId, "parentId")
		if err != nil {
			return err
		}

		data.ParentIdVal = &parentId
	} else {
		data.ParentIdVal = nil
	}
	log.Printf("pareniId %v", data.ParentIdVal)
	return nil
}

func (h *CategoryHandler) create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	token, ok := ctx.Value("token").(*service.UserToken)
	if !ok || token == nil {
		render.Render(w, r, ErrNotFound)
		return
	}

	data := &CategoryCreateRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	createRequest := &service.CategoryCreateRequest{
		Name:     data.Name,
		Currency: data.CurrencyVal,
		UserId:   token.UserId,
	}

	if data.ParentIdVal != nil {
		createRequest.ParentId = data.ParentIdVal
	}

	category, parent, err := h.categoryService.Create(context.Background(), createRequest)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	render.JSON(w, r, NewCategoryResponse(category, parent))
}
