// ./logset.go
package logger

import (
	"log"
	"os"
	"path/filepath"
	"time"
)

// initLogger инициализирует логирование, создает директорию для логов и файл для записи логов.
func InitLogger() {
	logDir := "./09_logs"

	// Проверяем, существует ли директория для логов
	CreateLogDirectory(logDir)

	// Создаем файл для логов с текущей датой и временем
	logFile := CreateLogFile(logDir)

	// Перенаправляем вывод логов в файл и устанавливаем формат логов
	SetupLogOutput(logFile)
}

// createLogDirectory проверяет существование директории для логов и создает ее при необходимости.
func CreateLogDirectory(logDir string) {
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		// Директория не существует, создаем ее
		if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
			log.Fatalf("Не удалось создать директорию для логов: %v", err)
		} else {
			log.Println("Директория для логов создана успешно.")
		}
	} else {
		// Директория уже существует
		log.Println("Директория для логов уже существует.")
	}
}

// createLogFile создает файл для логов с текущей датой и временем.
func CreateLogFile(logDir string) *os.File {
	logFileName := filepath.Join(logDir, time.Now().Format("2006-01-02_15-04-05")+".log")
	// os.O_CREATE: флаг указывает, что если файл не существует, он должен быть создан.
	// os.O_WRONLY: флаг указывает, что файл будет открыт только для записи. Это означает, что не получится читать данные из файла, только писать в него.
	// os.O_APPEND: флаг указывает, что данные, которые записываются в файл, должны добавляться в конец файла, а не перезаписывать его содержимое.
	logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Не удалось открыть или создать файл логов: %v", err)
	} else {
		log.Println("Файл логов успешно создан.")
	}
	return logFile
}

// setupLogOutput перенаправляет вывод логов в файл и устанавливает параметры логирования.
func SetupLogOutput(logFile *os.File) {
	// Перенаправляем вывод логов в файл и в консоль
	log.SetOutput(logFile)
	// Формат вывода: 2024/10/31 12:34:56 main.go:10: Это логовое сообщение
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}
