package repo

import (
	"context"
	"deez-nats/internal/models"
	"deez-nats/pkg/logging"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type permanentRepo interface {
	CreateOrder(ctx context.Context, order models.Order) error
	GetOrders(ctx context.Context) ([]models.Order, error)
}

type PermanentRepo struct {
	*pgxpool.Pool
	*logging.Logger
}

func (r *PermanentRepo) CreateOrder(ctx context.Context, order models.Order) error {
	tx, err := r.Begin(ctx)
	if err != nil {
		r.Error(err)
		return err
	}
	var rollbackError error
	defer func(tx pgx.Tx, ctx context.Context, err *error) {
		if *err != nil {
			r.Info("rollback transaction")
			rollbackErr := tx.Rollback(ctx)
			if rollbackErr != nil {
				r.Error("tx already closed")
			}
		}
	}(tx, ctx, &rollbackError)

	delivery := order.Delivery
	var deliveryID string
	err = tx.QueryRow(ctx, insertDeliveryQuery,
		delivery.Name,
		delivery.Phone,
		delivery.Zip,
		delivery.City,
		delivery.Address,
		delivery.Region,
		delivery.Email,
	).Scan(&deliveryID)
	if err != nil {
		r.Errorf("error inserting delivery %s to db: %v", delivery.Name, err)
		rollbackError = err
		return err
	}

	payment := order.Payment
	var paymentID string
	err = tx.QueryRow(ctx, insertPaymentQuery,
		payment.Transaction,
		payment.RequestID,
		payment.Currency,
		payment.Provider,
		payment.Amount,
		payment.PaymentDT,
		payment.Bank,
		payment.DeliveryCost,
		payment.GoodsTotal,
		payment.CustomFee,
	).Scan(&paymentID)
	if err != nil {
		r.Errorf("error inserting payment %s to db: %v", payment.Transaction, err)
		rollbackError = err
		return err
	}

	_, err = tx.Exec(ctx, insertOrderQuery,
		order.ID,
		order.TrackNumber,
		order.Entry,
		deliveryID,
		paymentID,
		order.Locale,
		order.InternalSignature,
		order.CustomerID,
		order.DeliveryService,
		order.ShardKey,
		order.SMID,
		order.DateCreated,
		order.OofShard,
	)
	if err != nil {
		r.Errorf("error while inserting order %s : %v", order.ID, err)
		rollbackError = err
		return err
	}

	items := order.Items
	for _, i := range items {
		_, err = tx.Exec(ctx, insertItemQuery,
			order.ID,
			i.ChartID,
			i.TrackNumber,
			i.Price,
			i.Rid,
			i.Name,
			i.Sale,
			i.Size,
			i.TotalPrice,
			i.NMID,
			i.Brand,
			i.Status,
		)
		if err != nil {
			r.Errorf("error while inserting item: %v", err)
			rollbackError = err
			return err
		}
	}

	if err = tx.Commit(ctx); err != nil {
		rollbackError = err
		return err
	}

	return nil
}

func (r *PermanentRepo) GetOrders(ctx context.Context) ([]models.Order, error) {
	var orders []models.Order

	rows, err := r.Query(ctx, selectOrdersQuery)
	if err != nil {
		r.Errorf("cannot find orders in db: %v", err)
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var order models.Order
		err = rows.Scan(&order.ID, &order.TrackNumber, &order.Entry, &order.Locale,
			&order.InternalSignature, &order.CustomerID, &order.DeliveryService, &order.ShardKey, &order.SMID,
			&order.DateCreated, &order.OofShard)
		if err != nil {
			r.Error(err)
			return nil, err
		}

		err = r.fetchOrder(ctx, &order)
		if err != nil {
			r.Error(err)
			return nil, err
		}

		orders = append(orders, order)
	}

	return orders, nil
}

func (r *PermanentRepo) fetchOrder(ctx context.Context, order *models.Order) error {
	row := r.QueryRow(ctx, selectDeliveryPaymentQuery, order.ID)
	var delivery models.Delivery
	var payment models.Payment

	err := row.Scan(&delivery.Name, &delivery.Phone, &delivery.Zip, &delivery.City, &delivery.Address, &delivery.Region,
		&delivery.Email, &payment.Transaction, &payment.RequestID, &payment.Currency, &payment.Provider,
		&payment.Amount, &payment.PaymentDT, &payment.Bank, &payment.DeliveryCost, &payment.GoodsTotal, &payment.CustomFee)
	if err != nil {
		r.Error(err)
		return err
	}

	items, err := r.getItems(ctx, order.ID)
	if err != nil {
		r.Error(err)
		return err
	}

	order.Delivery = delivery
	order.Payment = payment
	order.Items = items

	return nil
}

func (r *PermanentRepo) getItems(ctx context.Context, orderID string) ([]models.Item, error) {
	rows, err := r.Query(ctx, selectItemsQuery, orderID)
	if err != nil {
		r.Error(err)
		return nil, err
	}

	defer rows.Close()
	var items []models.Item
	for rows.Next() {
		var item models.Item
		err = rows.Scan(&item.ChartID, &item.TrackNumber, &item.Price, &item.Rid, &item.Name, &item.Sale, &item.Size,
			&item.TotalPrice, &item.NMID, &item.Brand, &item.Status)
		if err != nil {
			r.Error(err)
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}
