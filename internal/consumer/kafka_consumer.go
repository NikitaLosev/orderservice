// consumer: получение сообщений из Kafka и сохранение заказов
package consumer

import (
	"LZero/internal/service"
	"LZero/pkg/models"
	"context"
	"encoding/json"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/segmentio/kafka-go"
)

// StartKafkaConsumer инициализирует консьюмер Kafka и обрабатывает сообщения заказов
func StartKafkaConsumer(pool *pgxpool.Pool) {
	ctx := context.Background()
	// инициализация Kafka reader
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "orders_topic",
		GroupID: "orders_consumer",
	})
	defer reader.Close()

	log.Println("🔍 Консьюмер Kafka запущен")
	// цикл чтения сообщений из топика
	for {
		msg, err := reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("Ошибка чтения сообщения: %v", err)
			continue
		}

		// парсинг JSON в структуру Order
		var order models.Order
		if err := json.Unmarshal(msg.Value, &order); err != nil {
			log.Printf("Ошибка парсинга JSON: %v", err)
			continue
		}

		// сохранение заказа в БД и обновление кеша
		if err := service.SaveOrder(pool, order); err != nil {
			log.Printf("Ошибка сохранения заказа: %v", err)
			continue
		}
		// подтверждение обработки сообщения (commit)
		if err := reader.CommitMessages(ctx, msg); err != nil {
			log.Printf("Ошибка коммита сообщения: %v", err)
		}
		log.Printf("Заказ сохранён: %s", order.OrderUID)
	}
}
