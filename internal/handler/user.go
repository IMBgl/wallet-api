package handler

import (
	"context"
	"errors"
	"github.com/IMBgl/go-wallet-api/internal/domain"
	"github.com/IMBgl/go-wallet-api/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"log"
	"net/http"
)

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
	Email    string `json:"email" binding:"required,alphaunicode,email"`
	Password string `json:"password" binding:"required,min=5,max=100"`
}

type UserSignUpRequest struct {
	Name     string      `json:"name" validate:"required,ascii,max=25,min=5"`
	Email    interface{} `json:"email" validate:"required,email"`
	Password interface{} `json:"password" validate:"required,min=5,max=100"`
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

func (h *UserHandler) singUp(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	request := &UserSignUpRequest{}

	err := unmarshallRequest(r, request)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	err, verrs := validateRequest(request)
	if err != nil {
		log.Printf("validation process err %v", err)
		render.Render(w, r, ErrInvalidRequest(errors.New("invalid request")))
		return
	}

	if len(verrs) > 0 {
		render.JSON(w, r, verrs)
		return
	}
	serviceRequest := service.SignUpRequest{
		Name:     request.Name,
		Email:    request.Email.(string),
		Password: request.Password.(string),
	}

	us, token, err := h.userService.SingUp(ctx, serviceRequest)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	render.JSON(w, r, NewCredentialsResponse(us, token))
}

func (h *UserHandler) signIn(w http.ResponseWriter, r *http.Request) {
	//request := &UserSignInRequest{}
	//if err := render.Bind(r, request); err != nil {
	//	render.Render(w, r, ErrInvalidRequest(err))
	//	return
	//}
	//
	//dto := service.SignInRequest{
	//	Email:    request.Email,
	//	Password: request.Password,
	//}
	//
	//us, token, err := h.userService.SingIn(context.Background(), dto)
	//if err != nil {
	//	render.Render(w, r, ErrInvalidRequest(err))
	//	return
	//}
	//
	//render.JSON(w, r, NewCredentialsResponse(us, token))
}
