package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	// env "github.com/tricoderu/utils/env"
	// env "utils/env"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	*http.Server
	StopChan chan os.Signal
}

const (
	EnvLink string = "03_env/ports.env"
	WebVar  string = "TODO_FRONTEND_DIR"
)

// RunServer запускает HTTP сервер.
func RunServer() {

	// Загружаем переменные окружения из .env файла
	webDir := env.checkEnvVar(EnvLink, WebVar)

	// Создаем маршрутизатор
	r := chi.NewRouter()

	// Создаем обработчик статических файлов (то есть нашего frontend-а)
	r.Handle("/*", http.FileServer(http.Dir(webDir)))

	// Получаем номер порта
	port7540 := os.Getenv("TODO_PORT_7540")
	if port7540 == "" {
		log.Fatal("Фатальная ошибка: Переменная TODO_PORT_7540 не задана.")
	}

	// Получаем номер портов, на которых будут запущен(ы) сервер(ы)
	ports := []string{port7540}

	// Создаем WaitGroup для ожидания завершения работы серверной горутины
	var wg sync.WaitGroup

	// Массив для хранения ссылок на созданные сервера
	var servers []*Server

	// Создаем HTTP сервер, пока со всевозможными параметрами
	// Учтем возможность создания нескольких серверов для разных целей (для разных обработчиков)
	for _, port := range ports {
		srv := &Server{
			Server: &http.Server{
				Addr:         fmt.Sprintf(":%s", port),
				Handler:      r,
				ReadTimeout:  10 * time.Second,
				WriteTimeout: 10 * time.Second,
				IdleTimeout:  120 * time.Second,
			},
			StopChan: make(chan os.Signal, 1),
		}
		wg.Add(1)

		// Сохраняем ссылку на созданный сервер
		servers = append(servers, srv)

		// Запускаем каждый сервер в отдельной горутине
		go func(port string, srv *Server) {
			defer wg.Done()
			defer close(srv.StopChan)

			log.Printf("Запуск сервера на порту %s\n", port)

			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("Фатальная ошибка при запуске сервера: %v", err)
			} else {
				log.Printf("Сервер успешно запущен на порту %s\n", port)
			}
		}(port, srv)

		// Создаем канал, чтобы слушать и перенаправлять сигналы завершения из операционной системы (SIGINT и SIGTERM)
		// Прежде всего, чтобы можно было нормально закрыть сессию при нажатии Ctrl+C
		signal.Notify(srv.StopChan, syscall.SIGINT, syscall.SIGTERM)
	}

	log.Println("Ожидание остановки серверов...")

	// Запускаем обработку сигналов для каждого сервера
	for _, srv := range servers {
		ServerWithCancel(srv, &wg)
	}

	// Ждем завершения работы горутины сервера
	wg.Wait()
	log.Println("Все серверы остановлены.")
}
