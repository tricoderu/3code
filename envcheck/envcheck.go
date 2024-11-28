package env

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func CheckEnvVar(filePath, key string) (string, error) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return "", fmt.Errorf("директория с файлом %s не существует: %v", filePath, err)
	}

	err := godotenv.Load(filePath)
	if err != nil {
		return "", fmt.Errorf("не удалось загрузить файл .env: %v", err)
	}
	log.Println("Файл .env успешно загружен.")

	value := os.Getenv(key)
	if value == "" {
		return "", fmt.Errorf("переменная окружения %s не задана", key)
	}
	return value, nil
}
