package server

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"
)

// ServerTimeout завершает работу сервера с указанным тайм-аутом в секундах.
func ServerTimeout(ctx context.Context, srv *http.Server, seconds int) {
	// Устанавливаем таймаут для завершения работы сервера
	ctx, cancel := context.WithTimeout(ctx, time.Duration(seconds)*time.Second)
	// используем defer cancel() для гарантированного освобождения ресурсов, связанных с контекстом
	defer cancel()

	// Корректное завершение работы сервера
	log.Printf("Начато завершение работы сервера с таймаутом в %d секунд.\n", seconds)
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Фатальная ошибка при завершении работы сервера: %v", err)
	}
	log.Println("Сервер успешно остановлен.")
}

// ServerWithCancel завершает работу сервера при силовом обрыве.
func ServerWithCancel(srv *Server, wg *sync.WaitGroup) {
	<-srv.StopChan
	// Когда сервер получит сигнал из канала stop, то он начинает остановку и выводится сообщение о завершении работы сервера
	log.Printf("Получен сигнал из StopChan: %v. Начато завершение работы сервера...\n", srv.StopChan)

	// Завершаем работу сервера с таймаутом
	ServerTimeout(context.Background(), srv.Server, 5)

	wg.Done()
}

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
