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

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

type Server struct {
	// Встраиваем *http.Server в структуру Server через указатель на стандартный HTTP-сервер из пакета net/http
	*http.Server
	// Объявляем канал StopChan, который будет использоваться для передачи сигналов ОС
	StopChan chan os.Signal
}

const (
	envLink   string        = "02_env/ports.env"
	webVar    string        = "TODO_FRONTEND_DIR"
	port7540  string        = "TODO_PORT_7540"
	readTime  time.Duration = 10
	writeTime time.Duration = 10
	idleTime  time.Duration = 120
)

func RunServer() {
	// Загружаем конфигурацию
	if err := loadEnvConfig(); err != nil {
		log.Fatalf("Ошибка при загрузке конфигурации сервера: %v", err)
	}

	// Создаем экземпляр сервера
	srv := createServer()

	// Добавляем обработку сигналов
	handleSignals(srv)

	// Создаем WaitGroup для ожидания завершения работы серверной горутины
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		startServer(srv)
	}()

	// Обработка остановки сервера
	go ServerWithCancel(srv, &wg)

	log.Println("Ожидание остановки сервера...")
	wg.Wait()
	log.Println("Сервер остановлен.")
}

// Загружает настройки сервера из файла .env и возвращает ошибку, если не удалось.
func loadEnvConfig() error {
	if err := godotenv.Load(envLink); err != nil {
		return err
	}
	log.Println("Файл .env с конфигурацией сервера успешно загружен.")
	return nil
}

// Создаёт экземпляр сервера, настраивает маршруты для обслуживания статики и задаёт параметры подключения (порт и таймауты)
func createServer() *Server {
	webDir := os.Getenv(webVar)
	if webDir == "" {
		log.Fatalf("Фатальная ошибка: Переменная %v не задана.", webVar)
	}

	r := chi.NewRouter()
	r.Handle("/*", http.FileServer(http.Dir(webDir)))

	serverPort := os.Getenv(port7540)
	if serverPort == "" {
		log.Fatalf("Фатальная ошибка: Переменная %v не задана.", port7540)
	}

	srv := &Server{
		Server: &http.Server{
			Addr:         fmt.Sprintf(":%s", port7540),
			Handler:      r,
			ReadTimeout:  readTime * time.Second,
			WriteTimeout: writeTime * time.Second,
			IdleTimeout:  idleTime * time.Second,
		},
		StopChan: make(chan os.Signal, 1),
	}

	return srv
}

// Настраивает обработку сигналов ОС (SIGINT и SIGTERM), чтобы сервер мог корректно завершить свою работу.
func handleSignals(srv *Server) {
	signal.Notify(srv.StopChan, syscall.SIGINT, syscall.SIGTERM)
}

// Запускает HTTP сервер и логирует информацию о запуске.
func startServer(srv *Server) {
	log.Printf("Запуск сервера на %s\n", srv.Addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Фатальная ошибка при запуске сервера: %v", err)
	}
	log.Printf("Сервер успешно запущен на %s\n", srv.Addr)
}
