package repo

const (
	insertDeliveryQuery = `INSERT INTO deliveries(name, phone, zip, city, address, region, email) 
	VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	insertPaymentQuery = `INSERT INTO payments(transaction, request_id, currency, provider, amount, payment_dt, bank,
                     delivery_cost, goods_total, custom_fee) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id`
	insertItemQuery = `INSERT INTO items(order_id, chrt_id, track_number, price, rid, name, sale, size, total_price,
                  nm_id, brand, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
	insertOrderQuery = `INSERT INTO orders(id,track_number, entry, delivery_id, payment_id, locale,
                   internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard) VALUES   
                	($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`

	selectOrdersQuery = `SELECT 
    		id,
    		track_number,
    		entry,
    		locale,
    		internal_signature,
    		customer_id,
    		delivery_service,
    		shardkey,
    		sm_id,
    		date_created,
    		oof_shard
    	FROM orders`
	selectDeliveryPaymentQuery = `SELECT 
			d.name as delivery_name,
			d.phone as delivery_phone,
			d.zip as delivery_zip,
			d.city as delivery_city,
			d.address as delivery_address,
			d.region as delivery_region,
			d.email as delivery_email,
			p.transaction,
			p.request_id,
			p.currency,
			p.provider,
			p.amount,
			p.payment_dt,
			p.bank,
			p.delivery_cost,
			p.goods_total,
			p.custom_fee
		FROM orders o
		JOIN deliveries d ON o.delivery_id = d.id
		JOIN payments p ON o.payment_id = p.id
		WHERE o.id = $1`
	selectItemsQuery = `SELECT 
			i.chrt_id,
			i.track_number,
			i.price,
			i.rid,
			i.name,
			i.sale,
			i.size,
			i.total_price,
			i.nm_id,
			i.brand,
			i.status
		FROM items i
		WHERE i.order_id = $1`
)
