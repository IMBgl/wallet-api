package handler

import (
	"context"
	"github.com/IMBgl/go-wallet-api/internal/service"
	"github.com/go-chi/render"
	"log"
	"net/http"
)

const AUTH_HEADER = "X-Api-Key"

type apiMiddleware struct {
	service service.Service
}

func NewApiMiddleware(service service.Service) *apiMiddleware {
	return &apiMiddleware{service: service}
}

func (m *apiMiddleware) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenHeader := r.Header.Get(AUTH_HEADER)
		if tokenHeader == "" {
			log.Printf("invalid token header %s", tokenHeader)
			render.Render(w, r, ErrNotFound)
			return
		}

		token, err := m.service.Token().GetByValue(context.Background(), tokenHeader)
		if err != nil || token == nil {
			log.Printf("GetToken by value err %v", err)
			render.Render(w, r, ErrNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), "token", token)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
