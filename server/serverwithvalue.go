package server

import (
	"context"
	"log"
	"net/http"
	"time"
)

// ServerWithValue завершает работу сервера с контекстом, содержащим значения.
func ServerWithValue(srv *http.Server, key string, value interface{}, seconds int) {
	ctx, cancel := context.WithTimeout(context.WithValue(context.Background(), key, value), time.Duration(seconds)*time.Second)
	defer cancel()

	// Корректное завершение работы сервера
	log.Printf("Начато завершение работы сервера с контекстом, содержащим %v: %v, таймаут в %d секунд.\n", key, value, seconds)
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Фатальная ошибка при завершении работы сервера: %v", err)
	}
	log.Println("Сервер успешно остановлен.")
}
