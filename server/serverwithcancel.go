package server

import (
	"context"
	"log"
	"sync"
	"time"
)

// ServerWithCancel завершает работу сервера при силовом обрыве.
func ServerWithCancel(srv *Server, wg *sync.WaitGroup, timeout time.Duration) {
	<-srv.StopChan
	// Когда сервер получит сигнал из канала stop, то он начинает остановку и выводится сообщение о завершении работы сервера
	log.Println("Получен сигнал остановки работы сервера.")
	log.Println("Начато завершение работы сервера...")

	// Создаем новый контекст для таймаута при завершении работы сервера
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel() // гарантируем освобождение ресурсов

	// Завершаем работу сервера с таймаутом
	ServerTimeout(ctx, srv.Server)

	wg.Done()
}
