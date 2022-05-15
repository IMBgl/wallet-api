package repository

import (
	"context"
	"github.com/IMBgl/go-wallet-api/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"gorm.io/gorm"
	"time"
)

type CategoryModel struct {
	gorm.Model
	Id       uuid.UUID
	Name     string
	UserId   uuid.UUID
	Currency CurrencyValue
}

func (CategoryModel) TableName() string {
	return "categories"
}

func (m *CategoryModel) Entity() (*domain.Category, error) {
	return &domain.Category{
		Id:        m.Id,
		Name:      m.Name,
		UserId:    m.UserId,
		Currency:  m.Currency.Currency,
		CreatedAt: m.CreatedAt,
	}, nil
}

func (m *CategoryModel) FromEntity(e *domain.Category) *CategoryModel {
	m.Id = e.Id
	m.Name = e.Name
	m.UserId = e.UserId
	m.Currency = CurrencyValue{e.Currency}

	return m
}

type categoryRepository struct {
	repository
}

func CategoryRepository(conn *pgx.Conn) *categoryRepository {
	return &categoryRepository{repository{Conn: conn}}
}

func (r *categoryRepository) Save(ctx context.Context, c *domain.Category) error {
	_, err := r.Conn.Exec(ctx, `
				insert into categories (id, "name", user_id, parent_id, currency, created_at, updated_at)
													values($1,$2,$3,$4,$5,$6,$7)
													on conflict (id) do update 
													set name = $2, updated_at = $7;`,
		c.Id, c.Name, c.UserId, c.ParentId, c.Currency.Val(), c.CreatedAt, time.Now())

	return err
}

func (r *categoryRepository) GetById(ctx context.Context, id uuid.UUID) (*domain.Category, error) {
	category := domain.Category{}
	currencyVal := ""

	err := r.Conn.QueryRow(ctx, "select id, \"name\", user_id, parent_id, currency, created_at from categories where id=$1", id).Scan(&category.Id, &category.Name, &category.UserId, &category.ParentId, &currencyVal, &category.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	currency, err := domain.CurrencyFromString(currencyVal)
	if err != nil {
		return nil, err
	}

	category.Currency = currency

	return &category, nil
}

func (r *categoryRepository) FindByUserId(ctx context.Context, userId uuid.UUID) (list []*domain.Category, err error) {
	rows, _ := r.Conn.Query(ctx, "select id, \"name\", user_id, parent_id, currency, created_at from categories where user_id=$1", userId)

	for rows.Next() {
		i := domain.Category{}
		currencyVal := ""

		err := rows.Scan(&i.Id, &i.Name, &i.UserId, &i.ParentId, &currencyVal, &i.CreatedAt)
		if err != nil {
			return nil, err
		}

		currency, err := domain.CurrencyFromString(currencyVal)
		if err != nil {
			return nil, err
		}

		i.Currency = currency
		list = append(list, &i)
	}

	return
}

func (r *categoryRepository) FindByUserIdWithNullParent(ctx context.Context, userId uuid.UUID) (list []*domain.Category, err error) {
	rows, _ := r.Conn.Query(ctx, "select id, \"name\", user_id, parent_id, currency, created_at from categories where user_id=$1 and parent_id is null", userId)

	for rows.Next() {
		i := domain.Category{}
		currencyVal := ""

		err := rows.Scan(&i.Id, &i.Name, &i.UserId, &i.ParentId, &currencyVal, &i.CreatedAt)
		if err != nil {
			return nil, err
		}

		currency, err := domain.CurrencyFromString(currencyVal)
		if err != nil {
			return nil, err
		}

		i.Currency = currency
		list = append(list, &i)
	}

	return
}

func (r *categoryRepository) FindByIdAndUserId(ctx context.Context, id, userId uuid.UUID) (*domain.Category, error) {
	currencyVal := ""
	category := domain.Category{}
	err := r.Conn.QueryRow(ctx, "select id, \"name\", user_id, parent_id, currency, created_at from categories where id=$1 and user_id=$2", id, userId).Scan(&category.Id, &category.Name, &category.UserId, &category.ParentId, &currencyVal, &category.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	currency, err := domain.CurrencyFromString(currencyVal)
	if err != nil {
		return nil, err
	}

	category.Currency = currency

	return &category, nil
}

func (r *categoryRepository) Delete(ctx context.Context, c *domain.Category) error {
	_, err := r.Conn.Exec(ctx, "delete from categories where id=$1", c.Id)

	return err
}

func (r *categoryRepository) GetChildren(ctx context.Context, c *domain.Category) (list []*domain.Category, err error) {
	rows, _ := r.Conn.Query(ctx, "select id, \"name\", user_id, parent_id, currency, created_at from categories where parent_id=$1", c.Id)

	for rows.Next() {
		i := domain.Category{}
		currencyVal := ""

		err := rows.Scan(&i.Id, &i.Name, &i.UserId, &i.ParentId, &currencyVal, &i.CreatedAt)
		if err != nil {
			return nil, err
		}

		currency, err := domain.CurrencyFromString(currencyVal)
		if err != nil {
			return nil, err
		}

		i.Currency = currency
		list = append(list, &i)
	}

	return
}
