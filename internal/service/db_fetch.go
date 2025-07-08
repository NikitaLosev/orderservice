// service: выборка заказов из БД при cache-miss
package service

import (
	"LZero/pkg/models"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// GetOrderFromDB возвращает заказ из БД, false если не найден, или ошибку.
func GetOrderFromDB(pool *pgxpool.Pool, uid string) (models.Order, bool, error) {
	var o models.Order

	// выбор основных полей заказа
	err := pool.QueryRow(context.Background(), `
		SELECT order_uid, track_number, entry, locale, internal_signature,
		       customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard
		FROM orders WHERE order_uid=$1`, uid).Scan(
		&o.OrderUID, &o.TrackNumber, &o.Entry, &o.Locale, &o.InternalSignature,
		&o.CustomerID, &o.DeliveryService, &o.ShardKey, &o.SmID,
		&o.DateCreated, &o.OofShard,
	)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return models.Order{}, false, nil
		}
		return models.Order{}, false, fmt.Errorf("orders select: %w", err)
	}

	// загрузка информации о доставке
	_ = pool.QueryRow(context.Background(), `
		SELECT name, phone, zip, city, address, region, email
		FROM deliveries WHERE order_uid=$1`, uid).Scan(
		&o.Delivery.Name, &o.Delivery.Phone, &o.Delivery.Zip, &o.Delivery.City,
		&o.Delivery.Address, &o.Delivery.Region, &o.Delivery.Email,
	)

	// загрузка информации об оплате
	_ = pool.QueryRow(context.Background(), `
		SELECT transaction_id, request_id, currency, provider, amount, payment_dt,
		       bank, delivery_cost, goods_total, custom_fee
		FROM payments WHERE order_uid=$1`, uid).Scan(
		&o.Payment.Transaction, &o.Payment.RequestID, &o.Payment.Currency, &o.Payment.Provider,
		&o.Payment.Amount, &o.Payment.PaymentDT, &o.Payment.Bank, &o.Payment.DeliveryCost,
		&o.Payment.GoodsTotal, &o.Payment.CustomFee,
	)

	// загрузка позиций заказа
	rows, _ := pool.Query(context.Background(), `
		SELECT chrt_id, track_number, price, rid, name, sale, size,
		       total_price, nm_id, brand, status
		FROM items WHERE order_uid=$1`, uid)
	defer rows.Close()
	for rows.Next() {
		var it models.Item
		_ = rows.Scan(&it.ChrtID, &it.TrackNumber, &it.Price, &it.Rid, &it.Name,
			&it.Sale, &it.Size, &it.TotalPrice, &it.NmID, &it.Brand, &it.Status)
		o.Items = append(o.Items, it)
	}

	return o, true, nil
}
