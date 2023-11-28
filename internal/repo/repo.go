package repo

import (
	"context"
	"deez-nats/internal/models"
	"deez-nats/pkg/logging"
	"deez-nats/pkg/storage/cache"
	"github.com/jackc/pgx/v5/pgxpool"
)

type IRepository interface {
	AddOrder(ctx context.Context, order models.Order) error
	GetOrderById(id string) (models.Order, error)
	UploadCache(ctx context.Context) error
}

type Repository struct {
	permanentRepo
	CachedRepo
	l *logging.Logger
}

func NewRepository(db *pgxpool.Pool, cache cache.Cache, l logging.Logger) *Repository {
	return &Repository{
		permanentRepo: &PermanentRepo{db, &l},
		CachedRepo:    CachedRepo{&cache},
		l:             &l,
	}
}

func (r *Repository) AddOrder(ctx context.Context, order models.Order) error {
	err := r.CreateOrder(ctx, order)
	if err != nil {
		return err
	}

	err = r.InsertData(order)
	if err != nil {
		r.l.Errorf("error inserting data to cache %v", err)
		return err
	}
	r.l.Infof("order with id %s was successfully added to cache", order.ID)

	return nil
}

func (r *Repository) GetOrderById(id string) (models.Order, error) {
	return r.GetDataById(id)
}

func (r *Repository) UploadCache(ctx context.Context) error {
	orders, err := r.GetOrders(ctx)
	if err != nil {
		return err
	}

	for _, order := range orders {
		if _, ok := r.Cache.Data[order.ID]; !ok {
			err = r.InsertData(order)
			if err != nil {
				r.l.Errorf("error inserting data to cache %v", err)
				return err
			}
			r.l.Infof("order with id %s was successfully added to cache", order.ID)
		}
	}

	return nil
}
