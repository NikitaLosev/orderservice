// orders-producer читает JSON-файл и публикует его в Kafka (топик "orders_topic").
package main

import (
	"context"
	"flag"
	"io/ioutil"
	"log"

	"github.com/segmentio/kafka-go"
)

func main() {
	// разбор флагов командной строки
	file := flag.String("f", "test.json", "path to JSON file")
	flag.Parse()

	// чтение JSON из указанного файла
	data, err := ioutil.ReadFile(*file)
	if err != nil {
		log.Fatalf("read file: %v", err)
	}

	// инициализация писателя Kafka (брокеры и топик)
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "orders_topic",
	})
	defer w.Close()

	// отправка сообщения в Kafka
	if err := w.WriteMessages(context.Background(),
		kafka.Message{Value: data},
	); err != nil {
		log.Fatalf("write msg: %v", err)
	}
	log.Println("сообщение отправлено в Kafka")
}
