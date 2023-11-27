package migrations

import (
	"database/sql"
	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upCreateTables, downCreateTables)
}

func upCreateTables(tx *sql.Tx) error {
	_, err := tx.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
	create table if not exists deliveries(
	    id UUID PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		phone VARCHAR(20) NOT NULL,
		zip VARCHAR(10) NOT NULL,
		city VARCHAR(255) NOT NULL,
		address VARCHAR(255) NOT NULL,
		region VARCHAR(255) NOT NULL,
		email VARCHAR(255) NOT NULL
	);
	`)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
	create table if not exists payments (
	    id UUID PRIMARY KEY,
		transaction VARCHAR(255) NOT NULL,
		request_id VARCHAR(255) NOT NULL,
		currency VARCHAR(3) NOT NULL,
		provider VARCHAR(255) NOT NULL,
		amount BIGINT NOT NULL,
		payment_dt BIGINT NOT NULL,
		bank VARCHAR(255) NOT NULL,
		delivery_cost BIGINT NOT NULL,
		goods_total BIGINT NOT NULL,
		custom_fee BIGINT NOT NULL
	);
`)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
	CREATE TABLE IF NOT EXISTS orders(
	    id UUID PRIMARY KEY,
		order_uid VARCHAR(255) NOT NULL,
		track_number VARCHAR(255) NOT NULL,
		entry VARCHAR(255) NOT NULL,
		delivery_id UUID REFERENCES deliveries(id),
		payment_id UUID REFERENCES payments(id),
		locale VARCHAR(10) NOT NULL,
		internal_signature VARCHAR(255),
		customer_id VARCHAR(255) NOT NULL,
		delivery_service VARCHAR(255) NOT NULL,
		shardkey VARCHAR(255) NOT NULL,
		sm_id BIGINT NOT NULL,
		date_created TIMESTAMP WITH TIME ZONE NOT NULL,
		oof_shard VARCHAR(255) NOT NULL
	);`)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
	create table if not exists items(
	    id UUID PRIMARY KEY,
		order_id UUID REFERENCES orders(id),
		chrt_id BIGINT NOT NULL,
		track_number VARCHAR(255) NOT NULL,
		price BIGINT NOT NULL,
		rid VARCHAR(255) NOT NULL,
		name VARCHAR(255) NOT NULL,
		sale BIGINT NOT NULL,
		size VARCHAR(10) NOT NULL,
		total_price BIGINT NOT NULL,
		nm_id BIGINT NOT NULL,
		brand VARCHAR(255) NOT NULL,
		status INT NOT NULL
	)
`)

	return nil
}

func downCreateTables(tx *sql.Tx) error {
	_, err := tx.Exec("drop table deliveries cascade ")
	if err != nil {
		return err
	}
	_, err = tx.Exec("drop table payments cascade ")
	if err != nil {
		return err
	}
	_, err = tx.Exec("drop table orders cascade ")
	if err != nil {
		return err
	}
	_, err = tx.Exec("drop table items cascade ")
	if err != nil {
		return err
	}

	return nil
}
