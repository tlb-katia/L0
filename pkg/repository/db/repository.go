package db

import (
	"L0/config"
	"L0/entities"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"
	"time"
)

type Repository struct {
	DB  *sql.DB
	Rdb *redis.Client
}

func NewRepository(db *sql.DB, Rdb *redis.Client) *Repository {
	return &Repository{
		DB:  db,
		Rdb: Rdb,
	}
}

func Init(cnf *config.Config) (*sql.DB, error) {
	const op = "Repository.Init"

	conStr := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		cnf.PGHost, cnf.PGPort, cnf.PGUser, cnf.PGName, cnf.PGPassword)

	db, err := sql.Open("postgres", conStr)
	if err != nil {
		return nil, fmt.Errorf("%s  %s", op, err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("%s %s", op, err)
	}

	return db, nil
}

func (r *Repository) CreateOrder(ctx context.Context, order *entities.Order) error {
	const op = "Repository.CreateOrder"

	_, err := r.DB.ExecContext(ctx, `
        INSERT INTO orders (order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
    `, order.OrderUID, order.TrackNumber, order.Entry, order.Locale, order.InternalSignature, order.CustomerID, order.DeliveryService, order.Shardkey, order.SmID, order.DateCreated, order.OOFShard)

	if err != nil {
		return fmt.Errorf("%s %s", op, err)
	}

	_, err = r.DB.ExecContext(ctx, `
        INSERT INTO delivery (order_uid, name, phone, zip, city, address, region, email)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    `, order.OrderUID, order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip, order.Delivery.City, order.Delivery.Address, order.Delivery.Region, order.Delivery.Email)

	if err != nil {
		return fmt.Errorf("%s %s", op, err)
	}

	_, err = r.DB.ExecContext(ctx, `
        INSERT INTO payment (order_uid, transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
    `, order.OrderUID, order.Payment.Transaction, order.Payment.RequestID, order.Payment.Currency, order.Payment.Provider, order.Payment.Amount, order.Payment.PaymentDT, order.Payment.Bank, order.Payment.DeliveryCost, order.Payment.GoodsTotal, order.Payment.CustomFee)

	if err != nil {
		return fmt.Errorf("%s %s", op, err)
	}

	for _, item := range order.Items {
		_, err = r.DB.ExecContext(ctx, `
            INSERT INTO items (order_uid, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status)
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
        `, order.OrderUID, item.ChrtID, item.TrackNumber, item.Price, item.Rid, item.Name, item.Sale, item.Size, item.TotalPrice, item.NmID, item.Brand, item.Status)

		if err != nil {
			return fmt.Errorf("%s %s", op, err)
		}
	}

	r.Rdb.Set(ctx, fmt.Sprintf("order:%s", order.OrderUID), order, 10*time.Minute)

	return nil
}

func (r *Repository) GetOrderById(con context.Context, orderUID string) (*entities.Order, error) {
	const op = "Repository.GetOrderById"

	var order entities.Order
	var delivery entities.Delivery
	var payment entities.Payment
	var item entities.Item

	cachedOrder, err := r.Rdb.Get(con, fmt.Sprintf("order:%s", orderUID)).Result()
	if err == nil {
		if err := json.Unmarshal([]byte(cachedOrder), &order); err == nil {
			return &order, nil
		}
	}

	query := `
    SELECT
        o.order_uid,
        o.track_number,
        o.entry,
        o.locale,
        o.internal_signature,
        o.customer_id,
        o.delivery_service,
        o.shardkey,
        o.sm_id,
        o.date_created,
        o.oof_shard,
        d.name AS delivery_name,
        d.phone AS delivery_phone,
        d.zip AS delivery_zip,
        d.city AS delivery_city,
        d.address AS delivery_address,
        d.region AS delivery_region,
        d.email AS delivery_email,
        p.transaction AS payment_transaction,
        p.request_id AS payment_request_id,
        p.currency AS payment_currency,
        p.provider AS payment_provider,
        p.amount AS payment_amount,
        p.payment_dt AS payment_payment_dt,
        p.bank AS payment_bank,
        p.delivery_cost AS payment_delivery_cost,
        p.goods_total AS payment_goods_total,
        p.custom_fee AS payment_custom_fee,
        i.chrt_id AS item_chrt_id,
        i.track_number AS item_track_number,
        i.price AS item_price,
        i.rid AS item_rid,
        i.name AS item_name,
        i.sale AS item_sale,
        i.size AS item_size,
        i.total_price AS item_total_price,
        i.nm_id AS item_nm_id,
        i.brand AS item_brand,
        i.status AS item_status
    FROM 
        orders o
    LEFT JOIN 
        delivery d ON o.order_uid = d.order_uid
    LEFT JOIN 
        payment p ON o.order_uid = p.order_uid
    LEFT JOIN 
        items i ON o.order_uid = i.order_uid
    WHERE 
        o.order_uid = $1;
    `

	row := r.DB.QueryRowContext(con, query, orderUID)

	err = row.Scan(
		&order.OrderUID,
		&order.TrackNumber,
		&order.Entry,
		&order.Locale,
		&order.InternalSignature,
		&order.CustomerID,
		&order.DeliveryService,
		&order.Shardkey,
		&order.SmID,
		&order.DateCreated,
		&order.OOFShard,
		&delivery.Name,
		&delivery.Phone,
		&delivery.Zip,
		&delivery.City,
		&delivery.Address,
		&delivery.Region,
		&delivery.Email,
		&payment.Transaction,
		&payment.RequestID,
		&payment.Currency,
		&payment.Provider,
		&payment.Amount,
		&payment.PaymentDT,
		&payment.Bank,
		&payment.DeliveryCost,
		&payment.GoodsTotal,
		&payment.CustomFee,
		&item.ChrtID,
		&item.TrackNumber,
		&item.Price,
		&item.Rid,
		&item.Name,
		&item.Sale,
		&item.Size,
		&item.TotalPrice,
		&item.NmID,
		&item.Brand,
		&item.Status,
	)
	if err != nil {
		return nil, fmt.Errorf("%s %s", op, err)
	}

	order.Delivery = delivery
	order.Payment = payment

	return &order, nil
}
