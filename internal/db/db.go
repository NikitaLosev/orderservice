// db: инициализация подключения к PostgreSQL
package db

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ConnectDB создаёт и возвращает пул соединений к базе данных по ENV-конфигу
func ConnectDB() (*pgxpool.Pool, error) {
	// формирование строки подключения из переменных окружения
	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	// создаём контекст с таймаутом для инициализации пула
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// создаём пул соединений
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}
	// проверяем доступность базы (ping)
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping DB: %w", err)
	}
	return pool, nil
}
