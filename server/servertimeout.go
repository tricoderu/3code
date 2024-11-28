package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

// ServerTimeout завершает работу сервера с указанным тайм-аутом в секундах.
func ServerTimeout(ctx context.Context, srv *http.Server) error {
	// Извлекаем время дедлайна и проверяем, установлен ли он (с помощью переменной ok)
	deadline, ok := ctx.Deadline()

	if ok {
		log.Printf("Начато завершение работы сервера с таймаутом до %s.\n", deadline)
	} else {
		log.Println("Начато завершение работы сервера без установленного таймаута.")
	}

	// Проходимся по завершению сервера
	if err := srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("функция ServerTimeout: фатальная ошибка при завершении работы сервера: %w", err)
	}

	log.Println("Сервер успешно остановлен.")
	return nil // Возвращаем nil, если завершение прошло успешно
}
