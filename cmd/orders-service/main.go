// orders-service: тест JSON, подключение к БД, восстановление кеша, запуск консьюмера и HTTP-сервера
package main

import (
	"LZero/internal/consumer"
	"LZero/internal/db"
	"LZero/internal/server"
	"LZero/internal/service"
	"LZero/pkg/models"
	"encoding/json"
	"log"
	"os"
)

func main() {
	// Тест парсинга JSON (необязательно)
	raw, _ := os.ReadFile("test.json")
	var o models.Order
	if err := json.Unmarshal(raw, &o); err == nil {
		log.Printf("Парсинг JSON успешен, пример order_uid: %s", o.OrderUID)
	}

	// Подключение к PostgreSQL
	pool, err := db.ConnectDB()
	if err != nil {
		log.Fatalf("DB: %v", err)
	}
	defer pool.Close()
	log.Println("Подключение к БД выполнено")

	// Восстановление кеша из БД
	if err := service.RestoreCache(pool); err != nil {
		log.Fatalf("RestoreCache: %v", err)
	}
	log.Printf("Кеш загружен: %d заказов", len(service.OrderCache))

	// Запуск Kafka-консьюмера (в горутине) и HTTP-сервера
	go consumer.StartKafkaConsumer(pool)
	server.StartHTTPServer(pool)
}
