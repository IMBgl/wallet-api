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
	"net/http"
)

type TransactionHandler struct {
	transactionService service.TransactionService
	middleware         *apiMiddleware
}

func (h TransactionHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Use(h.middleware.Auth)
	r.Post("/", h.create)

	return r
}

type TransactionCreateRequest struct {
	Comment       string                 `json:"comment,omitempty"`
	Currency      string                 `json:"currency"`
	Type          string                 `json:"type"`
	Amount        interface{}            `json:"amount"`
	CategoryId    string                 `json:"categoryId"`
	WalletId      string                 `json:"walletId"`
	AmountVal     float32                `json:"-"`
	WalletIdVal   uuid.UUID              `json:"-"`
	CategoryIdVal uuid.UUID              `json:"-"`
	TypeVal       domain.TransactionType `json:"-"`
	CurrencyVal   domain.Currency        `json:"-"`
}

type TransactionResponse struct {
	Id         string  `json:"id"`
	Comment    string  `json:"comment"`
	Currency   string  `json:"currency"`
	Type       string  `json:"type"`
	Amount     float32 `json:"amount"`
	CategoryId string  `json:"categoryId"`
	WalletId   string  `json:"walletId"`
	UserId     string  `json:"userId"`
	CreatedAt  string  `json:"createdAt"`
}

func NewTransactionResponse(e *domain.Transaction) *TransactionResponse {
	return &TransactionResponse{
		Id:         e.Id.String(),
		UserId:     e.UserId.String(),
		CategoryId: e.CategoryId.String(),
		WalletId:   e.WalletId.String(),
		Currency:   e.Currency.Val(),
		Type:       e.Type.Val(),
		Comment:    e.Comment,
		CreatedAt:  e.CreatedAt.Format(DateTimeFormat()),
		Amount:     e.Amount,
	}
}

func (data *TransactionCreateRequest) Bind(r *http.Request) error {
	currency, err := domain.CurrencyFromString(data.Currency)
	if err != nil {
		return errors.New("currency value must be one of 'rur', 'eur, 'usd'")
	}
	data.CurrencyVal = currency

	transactionType, err := domain.TransactionTypeFromString(data.Type)
	if err != nil {
		return errors.New("type value must be one of 'in', 'out")
	}
	data.TypeVal = transactionType

	amountVal, err := validator.Float32(data.Amount, "amount")
	if err != nil {
		return err
	}
	data.AmountVal = float32(amountVal)

	walletIdVal, err := validator.Uuid(data.WalletId, "walletId")
	if err != nil {
		return err
	}
	data.WalletIdVal = walletIdVal

	categoryIdVal, err := validator.Uuid(data.CategoryId, "categoryId")
	if err != nil {
		return err
	}
	data.CategoryIdVal = categoryIdVal

	return nil
}

func (h *TransactionHandler) create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	token, ok := ctx.Value("token").(*service.UserToken)
	if !ok {
		render.Render(w, r, ErrNotFound)
		return
	}

	data := &TransactionCreateRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	createRequest := &service.TransactionCreateRequest{
		Comment:         data.Comment,
		Currency:        data.CurrencyVal,
		UserId:          token.UserId,
		WalletId:        data.WalletIdVal,
		CategoryId:      data.CategoryIdVal,
		Amount:          data.AmountVal,
		TransactionType: data.TypeVal,
	}

	category, err := h.transactionService.Create(context.Background(), createRequest)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	render.JSON(w, r, NewTransactionResponse(category))
}
