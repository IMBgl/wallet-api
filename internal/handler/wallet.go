package handler

import (
	"context"
	"errors"
	"github.com/IMBgl/go-wallet-api/internal/domain"
	"github.com/IMBgl/go-wallet-api/internal/service"
	"github.com/IMBgl/go-wallet-api/pkg/validator"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"net/http"
)

type WalletHandler struct {
	walletService service.WalletService
	middleware    *apiMiddleware
}

func (h WalletHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Use(h.middleware.Auth)
	r.Post("/", h.create)
	r.Get("/", h.getList)

	r.Route("/{walletId}", func(r chi.Router) {
		r.Delete("/", h.delete)
		r.Put("/", h.update)
	})

	return r
}

type WalletCreateRequest struct {
	Name       string      `json:"name"`
	Balance    interface{} `json:"balance,string"`
	Currency   string      `json:"currency"`
	BalanceVal float32     `json:"-"`
}

type WalletUpdateRequest struct {
	Name string `json:"name"`
}

type WalletResponse struct {
	Id       string  `json:"id"`
	Name     string  `json:"name"`
	Balance  float32 `json:"balance"`
	Currency string  `json:"currency"`
	UserId   string  `json:"userId"`
}

func NewWalletListResponse(wl []*domain.Wallet) []*WalletResponse {
	wlr := []*WalletResponse{}
	for _, w := range wl {
		wlr = append(wlr, NewWalletResponse(w))
	}
	return wlr
}

func NewWalletResponse(w *domain.Wallet) *WalletResponse {
	return &WalletResponse{
		Id:       w.Id.String(),
		UserId:   w.UserId.String(),
		Currency: w.Currency.Val(),

		Name:    w.Name,
		Balance: w.Balance,
	}
}

func (data *WalletCreateRequest) Bind(r *http.Request) error {
	if data.Name == "" {
		return errors.New("name field required")
	}

	balanceVal, err := validator.Float32(data.Balance, "balance")
	if err != nil {
		return err
	}
	data.BalanceVal = balanceVal

	return nil
}

func (data *WalletUpdateRequest) Bind(r *http.Request) error {
	if data.Name == "" {
		return errors.New("name field required")
	}
	return nil
}

func (h *WalletHandler) create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	token, ok := ctx.Value("token").(*service.UserToken)
	if !ok {
		render.Render(w, r, ErrNotFound)
		return
	}

	data := &WalletCreateRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	createRequest := &service.WalletCreateRequest{
		Name:     data.Name,
		Balance:  data.BalanceVal,
		Currency: data.Currency,
		UserId:   token.UserId,
	}

	wallet, err := h.walletService.Create(context.Background(), createRequest)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	render.JSON(w, r, NewWalletResponse(wallet))
}

func (h *WalletHandler) delete(w http.ResponseWriter, r *http.Request) {
	token := retrieveTokenOrFail(w, r)
	walletId := retrieveUuidOrFail(w, r, "walletId")

	deleteRequest := &service.WalletDeleteRequest{
		UserId:   token.UserId,
		WalletId: walletId,
	}

	err := h.walletService.Delete(context.Background(), deleteRequest)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	render.JSON(w, r, map[string]string{})
}

func (h *WalletHandler) getList(w http.ResponseWriter, r *http.Request) {
	token := retrieveTokenOrFail(w, r)

	getListRequest := &service.WalletGetListRequest{
		UserId: token.UserId,
	}

	walletList, err := h.walletService.GetList(context.Background(), getListRequest)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	render.JSON(w, r, NewWalletListResponse(walletList))
}

func (h *WalletHandler) update(w http.ResponseWriter, r *http.Request) {
	token := retrieveTokenOrFail(w, r)
	walletId := retrieveUuidOrFail(w, r, "walletId")

	data := &WalletUpdateRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	updateRequest := &service.WalletUpdateRequest{
		Name:     data.Name,
		UserId:   token.UserId,
		WalletId: walletId,
	}

	wallet, err := h.walletService.Update(context.Background(), updateRequest)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	render.JSON(w, r, NewWalletResponse(wallet))
}
