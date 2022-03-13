package handler

import (
	"errors"
	"fmt"
	"github.com/IMBgl/go-wallet-api/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"strconv"

	"github.com/go-chi/render"
	"net/http"
)

type apiHandler struct {
	service service.Service
}

func ApiHandler(s service.Service) *apiHandler {
	return &apiHandler{service: s}
}

func (h *apiHandler) Routes() *chi.Mux {
	mv := NewApiMiddleware(h.service)
	userHandler := &UserHandler{userService: h.service.User()}
	walletHandler := &WalletHandler{walletService: h.service.Wallet(), middleware: mv}
	categoryHandler := &CategoryHandler{categoryService: h.service.Category(), middleware: mv}
	transactionHandler := &TransactionHandler{transactionService: h.service.Transaction(), middleware: mv}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			render.JSON(w, r, map[string]string{"status": "ok"})
		})

		r.Mount("/user", userHandler.Routes())
		r.Mount("/wallet", walletHandler.Routes())
		r.Mount("/category", categoryHandler.Routes())
		r.Mount("/transaction", transactionHandler.Routes())
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		render.Status(r, 404)
		render.JSON(w, r, map[string]string{"error": "not found"})
	})

	return r
}

func DateTimeFormat() string {
	return fmt.Sprintf("2006-01-02 15:04:05")
}

func validateFloat32(value interface{}, name string) (float32, error) {
	switch t := value.(type) {
	case string:
		floatVal, err := strconv.ParseFloat(t, 32)
		if err != nil {
			return 0.0, errors.New(fmt.Sprintf("%s value must be float", name))
		}
		return float32(floatVal), nil
	case float32:
		return t, nil
	case float64:
		return float32(t), nil
	}

	return 0.0, errors.New(fmt.Sprintf("%s value must be float", name))
}

func retrieveTokenOrFail(w http.ResponseWriter, r *http.Request) *service.UserToken {
	ctx := r.Context()
	token, ok := ctx.Value("token").(*service.UserToken)
	if !ok {
		render.Render(w, r, ErrNotFound)
		return nil
	}

	return token
}

func retrieveUuidOrFail(w http.ResponseWriter, r *http.Request, paramName string) (uuidVal uuid.UUID) {
	param := chi.URLParam(r, paramName)
	uuidVal, err := uuid.Parse(param)
	if err != nil {
		render.Render(w, r, ErrNotFound)
		return
	}
	return
}
