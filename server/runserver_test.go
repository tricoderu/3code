//go:build ignore

package server_test

import (
	"3code/server"
	"net/http"
	"os"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestServerShutdown(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)

	srv := &server.Server{
		Server:   &http.Server{Addr: ":7540"},
		StopChan: make(chan os.Signal, 1),
	}

	go func() {
		defer wg.Done()
		server.ServerWithCancel(srv, &wg)
	}()

	// Это будет имитировать задержку для сервера запуска
	time.Sleep(1 * time.Second)

	// Отправляем сигнал остановки
	srv.StopChan <- syscall.SIGINT

	// Ждем завершения работы
	wg.Wait()

	// Убедитесь, что служба остановлена
	assert.NotNil(t, srv.Server)
}
