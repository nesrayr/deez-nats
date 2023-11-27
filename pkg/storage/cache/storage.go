package cache

import "deez-nats/internal/models"

type Cache struct {
	Data map[string]models.Order
}

func NewCache() Cache {
	return Cache{Data: make(map[string]models.Order)}
}
