// server: запуск HTTP-сервера и обработка запросов
package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"LZero/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

// StartHTTPServer запускает REST API и раздачу статических файлов
func StartHTTPServer(pool *pgxpool.Pool) {
	// Статика: / → static/index.html
	http.Handle("/", http.FileServer(http.Dir("static")))

	// API: GET /order/{order_uid}
	http.HandleFunc("/order/", func(w http.ResponseWriter, r *http.Request) {
		// 1) Лог каждого входящего запроса
		log.Printf("HTTP %s %s", r.Method, r.URL.Path)

		// 2) Извлечение UID
		uid := strings.TrimSpace(strings.TrimPrefix(r.URL.Path, "/order/"))
		log.Printf("Ищем заказ UID=%q", uid)

		// 3) Пытаемся из кеша
		order, ok := service.OrderCache[uid]
		if !ok {
			// 4) Фолбэк: из БД
			var err error
			order, ok, err = service.GetOrderFromDB(pool, uid)
			if err != nil {
				log.Printf("DB fetch error for %q: %v", uid, err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
				return
			}
			if ok {
				service.OrderCache[uid] = order
				log.Printf("Заказ из БД закеширован: %q", uid)
			}
		}

		if !ok {
			log.Printf("Заказ не найден: %q", uid)
			http.Error(w, "order not found", http.StatusNotFound)
			return
		}

		// 5) Отдаём JSON
		log.Printf("Отдаём JSON для заказа: %q", uid)
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(order); err != nil {
			log.Printf("JSON encode error for %q: %v", uid, err)
			http.Error(w, "encode error", http.StatusInternalServerError)
		}
	})

	addr := ":8081"
	log.Printf("HTTP сервер слушает на %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
