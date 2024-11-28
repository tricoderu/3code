//go:build ignore

package server

import (
	"context"
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

// envLink string = "02_env/server.env"
const envLink string = "./02_env/server.env"

var (
	frontEnd   string
	serverPort string

	readTime       time.Duration
	writeTime      time.Duration
	idleTime       time.Duration
	contextTimeout time.Duration
	timeUnit       string
)

func init() {
	// Загружаем конфигурацию
	log.Println("Загрузка конфигурации из .env файла...")
	loadEnvConfig()
}

// Загружает настройки сервера из файла .env и возвращает ошибку, если не удалось.
func loadEnvConfig() {
	if err := godotenv.Load(envLink); err != nil {
		log.Printf("Ошибка при загрузке .env файла: %v", err)
		log.Fatal("Фатальная остановка при попытке загрузке конфигурации сервера из .env файла")
	}

	frontEnd = os.Getenv("TODO_FRONTEND_DIR")
	serverPort = os.Getenv("TODO_PORT_7540")
	log.Printf("Загруженные переменные: frontEnd=%s, serverPort=%s", frontEnd, serverPort)

	readTime = getDurationFromEnv("SERVER_READ_TIME", 10)
	writeTime = getDurationFromEnv("SERVER_WRITE_TIME", 10)
	idleTime = getDurationFromEnv("SERVER_IDLE_TIME", 120)
	contextTimeout = getDurationFromEnv("CTX_TIMEOUT", 5)
	timeUnit = os.Getenv("TIME_UNIT")

	log.Println("Файл .env с конфигурацией сервера успешно загружен.")
}

func getDurationFromEnv(varName string, defaultValue int) time.Duration {
	value := os.Getenv(varName)
	if value == "" {
		log.Printf("Переменная окружения %s не установлена. Используем значение по умолчанию: %d", varName, defaultValue)
		return time.Duration(defaultValue) * time.Second
	}

	duration, err := time.ParseDuration(value + timeUnit)
	if err != nil {
		log.Fatalf("Ошибка при парсинге переменной окружения %s: %v", varName, err)
	}

	return duration
}

func RunServer() {
	// Создаем экземпляр сервера
	srv := createServer()

	// Добавляем обработку сигналов
	handleSignals(srv)

	// Создаем WaitGroup для ожидания завершения работы серверной горутины
	// У нас тут две горутины: одна для запуска сервера и одна для обработки остановки, поэтому 2
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		log.Println("Запускаем сервер...")
		if err := startServer(srv); err != nil {
			log.Printf("Ошибка запуска сервера: %v", err)
		}
	}()

	// Обработка остановки сервера
	go func() {
		defer wg.Done()
		log.Println("Ожидание сигнала остановки сервера...")
		<-srv.StopChan
		log.Println("Получен сигнал остановки сервера")

		ctx, cancel := context.WithTimeout(context.Background(), contextTimeout*time.Second)
		defer cancel()

		if err := ServerTimeout(ctx, srv.Server); err != nil {
			log.Printf("Ошибка при завершении работы сервера: %v", err)
		}
	}()

	log.Println("Ожидание остановки сервера...")
	wg.Wait()
	log.Println("Сервер остановлен.")
}

// Создаёт экземпляр сервера, настраивает маршруты для обслуживания статики и задаёт параметры подключения (порт и таймауты)
func createServer() *Server {
	r := chi.NewRouter()
	r.Handle("/*", http.FileServer(http.Dir(frontEnd)))

	if serverPort == "" {
		log.Fatalf("Фатальная ошибка: Переменная %v не задана.", serverPort)
	}

	srv := &Server{
		Server: &http.Server{
			Addr:         fmt.Sprintf(":%s", serverPort),
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
func startServer(srv *Server) error {
	log.Printf("Запуск сервера на %s\n", srv.Addr)

	// Попытка запустить сервер
	err := srv.ListenAndServe()
	if err != nil {
		if err == http.ErrServerClosed {
			// Сервер корректно остановлен, можем просто вернуть nilку
			log.Printf("Сервер успешно остановлен на %s", srv.Addr)
			return nil
		}
		// Если произошла другая ошибка, логируем и возвращаем её
		return fmt.Errorf("функция startServer: фатальная ошибка при запуске сервера: %w", err)
	}

	// Этот лог не будет достигнут, так как ListenAndServe блокирует выполнение, поэтому он не нужен
	log.Printf("Сервер успешно запущен на %s\n", srv.Addr)
	return nil // Возвращаем nil, если сервер успешен
}
