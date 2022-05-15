package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"github.com/IMBgl/go-wallet-api/internal/domain"
	"github.com/google/uuid"
	"math/rand"
	"time"
)

const TOKEN_LENGHT = 50

type tokenServiсe struct {
	repo Repository
}

type UserToken struct {
	Id        uuid.UUID
	UserId    uuid.UUID
	Exp       time.Time
	Value     string
	CreatedAt time.Time
}

type TokenRepository interface {
	GetById(ctx context.Context, id uuid.UUID) (*UserToken, error)
	GetByValue(ctx context.Context, value string) (*UserToken, error)
	FindByUser(ctx context.Context, user *domain.User) ([]*UserToken, error)
	Save(ctx context.Context, token *UserToken) error
	Delete(ctx context.Context, token *UserToken) error
	DeleteAllForUserAndSave(ctx context.Context, user *domain.User, token *UserToken) error
}

func NewTokenService(r Repository) *tokenServiсe {
	return &tokenServiсe{repo: r}
}

func NewUserToken(value string, userId uuid.UUID) *UserToken {
	return &UserToken{
		Id:        uuid.New(),
		UserId:    userId,
		Exp:       time.Now().Add(time.Hour * 24),
		Value:     value,
		CreatedAt: time.Now(),
	}
}

func (s *tokenServiсe) CreateForUser(u *domain.User) *UserToken {
	return NewUserToken(GenerateTokenValue(), u.Id)
}

func GenerateTokenValue() string {
	h := sha256.New()
	b := make([]byte, TOKEN_LENGHT)
	rand.Read(b)

	hash := hex.EncodeToString(h.Sum(b))

	return hash
}

func (s *tokenServiсe) GetValidByValue(ctx context.Context, v string) (*UserToken, error) {
	token, err := s.repo.Token().GetByValue(ctx, v)
	if err != nil {
		return nil, err
	}

	if token == nil || !token.IsValid() {
		return nil, ErrNotFound
	}

	return token, nil
}

func (t *UserToken) IsValid() bool {
	return t.Exp.After(time.Now())
}
