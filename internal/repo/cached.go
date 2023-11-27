package repo

import (
	"deez-nats/internal/models"
	"deez-nats/pkg/storage/cache"
	"errors"
)

type CachedRepo struct {
	*cache.Cache
}

func (r *CachedRepo) GetDataById(id string) (models.Order, error) {
	order, ok := r.Data[id]
	if !ok {
		return models.Order{}, errors.New("no order with such id")
	}
	return order, nil
}

func (r *CachedRepo) InsertData(order models.Order) error {
	if _, ok := r.Data[order.ID]; ok {
		return errors.New("order with such id already exists")
	}
	r.Data[order.ID] = order
	return nil
}
