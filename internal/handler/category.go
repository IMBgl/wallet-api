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
	r.Get("/", h.getList)

	r.Route("/{categoryId}", func(r chi.Router) {
		r.Delete("/", h.delete)
		r.Put("/", h.update)
		r.Get("/", h.getOne)
	})

	return r
}

type CategoryCreateRequest struct {
	Name        string          `json:"name"`
	Currency    string          `json:"currency"`
	ParentId    *string         `json:"parentId,omitempty"`
	ParentIdVal *uuid.UUID      `json:"-"`
	CurrencyVal domain.Currency `json:"-"`
}

type CategoryUpdateRequest struct {
	Name string `json:"name"`
}

type CategoryResponse struct {
	Id        string     `json:"id"`
	Name      string     `json:"name"`
	Currency  string     `json:"currency"`
	UserId    string     `json:"userId"`
	CreatedAt string     `json:"createdAt"`
	ParentId  *uuid.UUID `json:"parentId"`
}

type CategoryNodeResponse struct {
	*CategoryResponse
	Parent   *CategoryResponse       `json:"parent"`
	Children []*CategoryNodeResponse `json:"children"`
}

func NewCategoryNodeResponse(n *service.CategoryTreeNode) *CategoryNodeResponse {
	response := &CategoryNodeResponse{}
	response.CategoryResponse = NewCategoryResponse(n.Category)
	response.Parent = NewCategoryResponse(n.Parent)

	for _, childNode := range n.Children {
		response.Children = append(response.Children, NewCategoryNodeResponse(childNode))
	}

	return response
}

func NewCategoryListResponse(cList []*domain.Category) []*CategoryResponse {
	var responseList []*CategoryResponse
	for _, c := range cList {
		responseList = append(responseList, NewCategoryResponse(c))
	}

	return responseList
}

func NewCategoryResponse(c *domain.Category) *CategoryResponse {
	if c == nil {
		return nil
	}

	resp := &CategoryResponse{
		Id:        c.Id.String(),
		Name:      c.Name,
		UserId:    c.UserId.String(),
		Currency:  c.Currency.Val(),
		CreatedAt: c.CreatedAt.Format(DateTimeFormat()),
		ParentId:  c.ParentId,
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

func (data *CategoryUpdateRequest) Bind(r *http.Request) error {
	if data.Name == "" {
		return errors.New("name field required")
	}
	return nil
}

func (h *CategoryHandler) create(w http.ResponseWriter, r *http.Request) {
	token := retrieveTokenOrFail(w, r)
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

	category, err := h.categoryService.Create(context.Background(), createRequest)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	render.JSON(w, r, NewCategoryResponse(category))
}

func (h *CategoryHandler) getList(w http.ResponseWriter, r *http.Request) {
	token := retrieveTokenOrFail(w, r)

	serviceRequest := &service.CategoryGetListRequest{
		UserId: token.UserId,
	}

	categoryList, err := h.categoryService.GetList(context.Background(), serviceRequest)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	render.JSON(w, r, NewCategoryListResponse(categoryList))
}

func (h *CategoryHandler) delete(w http.ResponseWriter, r *http.Request) {
	token := retrieveTokenOrFail(w, r)
	categoryId := retrieveUuidOrFail(w, r, "categoryId")

	serviceRequest := &service.CategoryDeleteRequest{
		UserId:     token.UserId,
		CategoryId: categoryId,
	}

	err := h.categoryService.Delete(context.Background(), serviceRequest)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	render.JSON(w, r, map[string]string{})
}

func (h *CategoryHandler) update(w http.ResponseWriter, r *http.Request) {
	token := retrieveTokenOrFail(w, r)
	categoryId := retrieveUuidOrFail(w, r, "categoryId")
	data := &CategoryUpdateRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	serviceRequest := &service.CategoryUpdateRequest{
		Name:       data.Name,
		UserId:     token.UserId,
		CategoryId: categoryId,
	}

	node, err := h.categoryService.Update(context.Background(), serviceRequest)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	render.JSON(w, r, NewCategoryNodeResponse(node))
}

func (h *CategoryHandler) getOne(w http.ResponseWriter, r *http.Request) {
	token := retrieveTokenOrFail(w, r)
	categoryId := retrieveUuidOrFail(w, r, "categoryId")

	serviceRequest := &service.CategoryGetOneRequest{
		UserId:     token.UserId,
		CategoryId: categoryId,
	}

	categoryNode, err := h.categoryService.GetOne(context.Background(), serviceRequest)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	render.JSON(w, r, NewCategoryNodeResponse(categoryNode))
}
