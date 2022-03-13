package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"github.com/IMBgl/go-wallet-api/internal/domain"
	"github.com/google/uuid"
	"log"
)

const MAX_USER_TOKENS_COUNT = 5

type userService struct {
	repo         Repository
	tokenService TokenService
}

type SignUpRequest struct {
	Name     string
	Email    string
	Password string
}

type SignInRequest struct {
	Email    string
	Password string
}

type UserRepository interface {
	GetById(ctx context.Context, id uuid.UUID) (*domain.User, error)
	Save(ctx context.Context, u *domain.User) error
	Delete(ctx context.Context, u *domain.User) error
	SaveUserWithToken(ctx context.Context, u *domain.User, t *UserToken) error
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
}

func NewUserService(r Repository, ts TokenService) *userService {
	return &userService{repo: r, tokenService: ts}
}

func HashPassword(pass string) string {
	h := sha256.New()
	h.Write([]byte(pass))
	hash := hex.EncodeToString(h.Sum(nil))

	return hash
}

func (s *userService) SingUp(ctx context.Context, signUp SignUpRequest) (*domain.User, *UserToken, error) {
	found, err := s.repo.User().GetByEmail(ctx, signUp.Email)
	if err != nil {
		return nil, nil, err
	}

	if found != nil {
		return nil, nil, ErrEmailAlreadyInUse
	}

	user := domain.NewUser(signUp.Name, signUp.Email, HashPassword(signUp.Password))
	token := s.tokenService.CreateForUser(user)

	err = s.repo.User().SaveUserWithToken(ctx, user, token)
	if err != nil {
		return nil, nil, err
	}

	return user, token, nil
}

func (s *userService) SingIn(ctx context.Context, request SignInRequest) (*domain.User, *UserToken, error) {
	user, err := s.repo.User().GetByEmail(ctx, request.Email)
	if err != nil {
		return nil, nil, err
	}

	if user == nil {
		return nil, nil, ErrUserNotFound
	}

	err = s.CheckUserPassword(user, request.Password)
	if err != nil {
		return nil, nil, err
	}

	token, err := s.GetTokenForUser(ctx, user)
	if err != nil {
		return nil, nil, err
	}

	return user, token, nil
}

func (s *userService) GetTokenForUser(ctx context.Context, user *domain.User) (*UserToken, error) {
	tokenList, err := s.repo.Token().FindByUser(ctx, user)
	if err != nil {
		return nil, err
	}

	token := s.tokenService.CreateForUser(user)
	log.Printf("deletion error %+v", tokenList)
	if len(tokenList) > MAX_USER_TOKENS_COUNT-1 {
		err = s.repo.Token().DeleteAllForUserAndSave(ctx, user, token)
		log.Printf("deletion error %v", err)
	} else {
		err = s.repo.Token().Save(ctx, token)
	}
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (s *userService) CheckUserPassword(u *domain.User, p string) error {
	if u.Password != HashPassword(p) {
		return ErrInvalidCredentials
	}

	return nil
}
