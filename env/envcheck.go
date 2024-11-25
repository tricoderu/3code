package env

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Функция проверяет наличие файла .env по определенному пути filePath и наличие переменной key внутри него
func checkEnvVar(filePath, key string) string {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Fatalf("Директория с файлом %s не существует: %v", filePath, err)
	}

	err := godotenv.Load(filePath)
	if err != nil {
		log.Fatalf("Фатальная ошибка: Не удалось загрузить файл .env: %v", err)
	} else {
		log.Println("Файл .env успешно загружен.")
	}

	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Фатальная ошибка: Переменная %s не задана.", key)
	}
	return value
}
