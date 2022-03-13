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

type UserModel struct {
	gorm.Model
	Id        uuid.UUID
	Name      string
	Password  string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (UserModel) TableName() string {
	return "users"
}

func (ur *UserModel) Entity() (*domain.User, error) {
	return &domain.User{
		ur.Id,
		ur.Name,
		ur.Email,
		ur.Password,
		ur.CreatedAt,
	}, nil
}

func (um *UserModel) FromEntity(u *domain.User) *UserModel {
	um.Id = u.Id
	um.Name = u.Name
	um.Email = u.Email
	um.Password = u.Password
	um.CreatedAt = u.CreatedAt
	um.UpdatedAt = time.Now()

	return um
}

type userRepository struct {
	repository
}

func UserRepository(conn *pgx.Conn) *userRepository {
	return &userRepository{repository{Conn: conn}}
}

func (r *userRepository) GetById(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	user := domain.User{}

	err := r.Conn.QueryRow(ctx, "select id, name, email, password, created_at from users where id=$1", id).Scan(&user.Id, &user.Name, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
	}

	return &user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	user := domain.User{}

	err := r.Conn.QueryRow(ctx, "select id, name, email, password, created_at from users where email=$1", email).Scan(&user.Id, &user.Name, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) Save(ctx context.Context, u *domain.User) error {
	_, err := r.Conn.Exec(ctx, "insert into users (id, name, email, password, created_at, updated_at) values($1,$2,$3,$4,$5,$6)", u.Id, u.Name, u.Email, u.Password, u.CreatedAt, time.Now())

	return err
}

func (r *userRepository) Delete(ctx context.Context, u *domain.User) error {
	_, err := r.Conn.Exec(ctx, "delete from users where id = $1", u.Id)

	return err
}

func (r *userRepository) SaveUserWithToken(ctx context.Context, u *domain.User, t *service.UserToken) error {
	tx, err := r.Conn.Begin(ctx)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, "insert into users (id, name, email, password, created_at, updated_at) values($1,$2,$3,$4,$5,$6)", u.Id, u.Name, u.Email, u.Password, u.CreatedAt, time.Now())
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
