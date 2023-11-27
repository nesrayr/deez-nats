package cache

import "deez-nats/internal/models"

type Cache struct {
	Data map[string]models.Order
}
