package repo

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type permanentRepo interface {
}

type PermanentRepo struct {
	*pgxpool.Pool
}

//func (r *PermanentRepo) CreateOrder(order models.Order) error {
//	tx := r.Create(&order)
//	if tx.Error != nil {
//		return tx.Error
//	}
//	return nil
//}
//
//func (r *PermanentRepo) GetOrders() ([]models.Order, error) {
//	var orders []models.Order
//	tx := r.Table("orders").Find(&orders)
//	if tx.Error != nil {
//		return nil, tx.Error
//	}
//	return orders, nil
//}
