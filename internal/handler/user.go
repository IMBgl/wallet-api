package handler

import (
	"context"
	"errors"
	"github.com/IMBgl/go-wallet-api/internal/domain"
	"github.com/IMBgl/go-wallet-api/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"net/http"
)

type UserSignUpRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserHandler struct {
	userService service.UserService
}

func (h UserHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/singUp", h.singUp)
	r.Post("/singIn", h.signIn)
	return r
}

type UserSignInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CredentialsResponse struct {
	UserId   string `json:"userId"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Token    string `json:"token"`
	TokenExp string `json:"tokenExp"`
}

func NewCredentialsResponse(u *domain.User, t *service.UserToken) *CredentialsResponse {
	return &CredentialsResponse{
		UserId:   u.Id.String(),
		Email:    u.Email,
		Name:     u.Name,
		Token:    t.Value,
		TokenExp: t.Exp.Format(DateTimeFormat()),
	}
}

func (u *UserSignUpRequest) Bind(r *http.Request) error {
	if u.Name == "" {
		return errors.New("name field required")
	}
	if u.Email == "" {
		return errors.New("email field required")
	}
	if u.Password == "" {
		return errors.New("password field required")
	}

	if len(u.Password) < 6 {
		return errors.New("password must be at least 6 characters long")
	}

	return nil
}

func (u *UserSignInRequest) Bind(r *http.Request) error {
	if u.Email == "" {
		return errors.New("email field required")
	}
	if u.Password == "" {
		return errors.New("password field required")
	}

	if len(u.Password) < 6 {
		return errors.New("password must be at least 6 characters long")
	}

	return nil
}

func (h *UserHandler) singUp(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	data := &UserSignUpRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	cud := service.SignUpRequest{
		Name:     data.Name,
		Email:    data.Email,
		Password: data.Password,
	}

	us, token, err := h.userService.SingUp(ctx, cud)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	render.JSON(w, r, NewCredentialsResponse(us, token))
}

func (h *UserHandler) signIn(w http.ResponseWriter, r *http.Request) {
	request := &UserSignInRequest{}
	if err := render.Bind(r, request); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	dto := service.SignInRequest{
		Email:    request.Email,
		Password: request.Password,
	}

	us, token, err := h.userService.SingIn(context.Background(), dto)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	render.JSON(w, r, NewCredentialsResponse(us, token))
}
