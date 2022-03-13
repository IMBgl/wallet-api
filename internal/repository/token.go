package repository

import (
	"context"
	"github.com/IMBgl/go-wallet-api/internal/domain"
	"github.com/IMBgl/go-wallet-api/internal/service"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"gorm.io/gorm"
	"time"
)

type TokenModel struct {
	gorm.Model
	Id        uuid.UUID
	UserId    uuid.UUID
	ExpiresAt time.Time
	Hash      string
}

func (TokenModel) TableName() string {
	return "user_tokens"
}

func (m *TokenModel) Entity() (*service.UserToken, error) {
	return &service.UserToken{
		Id:        m.Id,
		UserId:    m.UserId,
		Exp:       m.ExpiresAt,
		Value:     m.Hash,
		CreatedAt: m.CreatedAt,
	}, nil
}

func (m *TokenModel) FromEntity(u *service.UserToken) *TokenModel {
	m.Id = u.Id
	m.UserId = u.UserId
	m.ExpiresAt = u.Exp
	m.Hash = u.Value
	m.CreatedAt = u.CreatedAt
	m.UpdatedAt = time.Now()

	return m
}

func MapUserTokenList(modelLsit []TokenModel) ([]*service.UserToken, error) {
	var list []*service.UserToken
	for _, model := range modelLsit {
		entity, err := model.Entity()
		if err != nil {
			return list, err
		}
		list = append(list, entity)
	}

	return list, nil
}

type tokenRepository struct {
	repository
}

func TokenRepository(conn *pgx.Conn) *tokenRepository {
	return &tokenRepository{repository{Conn: conn}}
}

func (r *tokenRepository) GetById(ctx context.Context, id uuid.UUID) (*service.UserToken, error) {
	token := service.UserToken{}

	err := r.Conn.QueryRow(ctx, "select id, user_id, hash, expires_at, created_at, updated_at from user_tokens where id=$1", id).Scan(&token.Id, &token.UserId, &token.Value, &token.Exp, &token.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, service.ErrNotFound
		}
		return nil, err
	}

	return &token, nil
}

func (r *tokenRepository) GetByValue(ctx context.Context, value string) (*service.UserToken, error) {
	token := service.UserToken{}

	err := r.Conn.QueryRow(ctx, "select id, user_id, hash, expires_at, created_at from user_tokens where hash=$1", value).Scan(&token.Id, &token.UserId, &token.Value, &token.Exp, &token.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &token, nil
}

func (r *tokenRepository) FindByUser(ctx context.Context, user *domain.User) ([]*service.UserToken, error) {
	list := []*service.UserToken{}
	rows, _ := r.Conn.Query(ctx, "select id, user_id, hash, expires_at, created_at from user_tokens where user_id=$1", user.Id)

	for rows.Next() {
		token := service.UserToken{}
		err := rows.Scan(&token.Id, &token.UserId, &token.Value, &token.Exp, &token.CreatedAt)
		if err != nil {
			return nil, err
		}
		list = append(list, &token)
	}

	return list, nil
}

func (r *tokenRepository) Save(ctx context.Context, t *service.UserToken) error {
	_, err := r.Conn.Exec(ctx, "insert into user_tokens (id, user_id, hash, expires_at, created_at, updated_at) values($1,$2,$3,$4,$5,$6)", t.Id, t.UserId, t.Value, t.Exp, t.CreatedAt, time.Now())

	return err
}

func (r *tokenRepository) Delete(ctx context.Context, t *service.UserToken) error {
	_, err := r.Conn.Exec(ctx, "delete from user_tokens where id = $1", t.Id)

	return err
}

func (r *tokenRepository) DeleteAllForUserAndSave(ctx context.Context, u *domain.User, t *service.UserToken) error {
	tx, err := r.Conn.Begin(ctx)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, "delete from user_tokens where user_id = $1", u.Id)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	_, err = tx.Exec(ctx, "insert into user_tokens (id, user_id, hash, expires_at, created_at, updated_at) values($1,$2,$3,$4,$5,$6)", t.Id, t.UserId, t.Value, t.Exp, t.CreatedAt, time.Now())
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}
