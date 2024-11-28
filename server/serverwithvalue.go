package server

import (
	"context"
	"log"
	"net/http"
	"time"
)

// Определяем свой собственный тип для ключа контекста
type contextKey string

const (
	// Создаем ключ для контекста
	keyContextValue contextKey = "keyContextValue"
)

// ServerWithValue завершает работу сервера с контекстом, содержащим значения.
func ServerWithValue(srv *http.Server, key contextKey, value interface{}, seconds int) {
	ctx, cancel := context.WithTimeout(context.WithValue(context.Background(), keyContextValue, value), time.Duration(seconds)*time.Second)
	defer cancel()

	// Корректное завершение работы сервера
	log.Printf("Начато завершение работы сервера с контекстом, содержащим %v: %v, таймаут в %d секунд.\n", keyContextValue, value, seconds)
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Фатальная ошибка при завершении работы сервера: %v", err)
	}
	log.Println("Сервер успешно остановлен.")
}
