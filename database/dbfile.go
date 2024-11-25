package database

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// checkIfDBExists проверяет существование указанного файла базы данных.
func CheckIfDBExists(dbFile string) (bool, error) {
	// Получаем информацию о файле
	_, err := os.Stat(dbFile)

	// Обрабатываем ошибки + можно добавить еще какие-то
	switch {
	case os.IsNotExist(err):
		log.Printf("Файл базы данных отсутствует: %v", err)
		return false, nil // Файл не существует
	case os.IsPermission(err):
		log.Printf("Ошибка доступа: у вас нет прав для доступа к файлу базы данных %s", dbFile)
		return false, err // Отказ в доступе
	case err != nil:
		log.Printf("Другая ошибка при проверке существования файла базы данных: %v", err)
		return false, err // Возвращаем ошибку, если какая-то иная
	}
	log.Printf("Файл базы данных существует по адресу %v", dbFile)
	return true, nil // Файл существует
}

// createDBDirectory создает директорию для базы данных.
func CreateDBDirectory(dbFile string) {
	// Получаем директорию из пути к файлу базы данных
	dbDir := filepath.Dir(dbFile)

	// Проверяем существование директории
	if _, err := os.Stat(dbDir); os.IsNotExist(err) {
		// Если директория не существует, создаем ее
		if err := os.MkdirAll(dbDir, os.ModePerm); err != nil {
			log.Fatalf("Не удалось создать директорию для базы данных: %v", err)
		} else {
			log.Println("Директория для базы данных успешно создана:", dbDir)
		}
	} else {
		log.Printf("Директория для базы данных уже существует: %v", dbDir)
		log.Println("Переходим к созданию файла базы данных...")
	}
}

// createDatabaseFile создает файл базы данных.
func CreateDatabaseFile(dbFile string) {
	// Создаем файл базы данных
	file, err := os.Create(dbFile)
	if err != nil {
		log.Fatalf("Не удалось создать файл базы данных: %v", err)
	}
	defer file.Close() // Закрываем файл после завершения работы
	log.Println("Файл базы данных успешно создан:", dbFile)
}

// setupDatabase проверяет существование базы данных и создает её, если она не найдена.
func SetupDatabase() string {
	var err error

	// Загружаем переменные окружения из .env файла
	err = godotenv.Load("03_env/db.env")
	if err != nil {
		log.Fatalf("Фатальная ошибка: Не удалось загрузить файл .env: %v", err)
	} else {
		log.Println("Файл .env успешно загружен.")
	}

	// Получаем значение переменной окружения TODO_DBFILE
	dbFile := os.Getenv("TODO_DBFILE")

	// Если переменная окружения не установлена, задаем путь по умолчанию
	if dbFile == "" {
		dbFile = "./db/scheduler.db" // Задаем относительный путь для файла БД
		log.Println("Используется путь к файлу базы данных по умолчанию:", dbFile)
	} else {
		log.Printf("Используется путь к файлу базы данных из переменной окружения: %v\n", dbFile)
	}

	// Проверяем, существует ли файл базы данных
	var exists bool
	counter := 0

	// Получаем количество попыток из переменной окружения
	dbAttemptsStr := os.Getenv("TODO_ATTEMPTS")
	if dbAttemptsStr == "" {
		log.Fatal("Фатальная ошибка: Переменная TODO_ATTEMPTS не задана.")
	}

	dbAttemptsInt, err := strconv.Atoi(dbAttemptsStr)
	log.Printf("Значение количества попыток из файла .env %v", dbAttemptsInt)
	if err != nil || dbAttemptsInt <= 0 {
		log.Printf("Некорректное значение количества попыток. Будет установлено значение по умолчанию: 3")
		dbAttemptsInt = 3 // Устанавливаем значение по умолчанию
	}

	for attempts := 0; attempts < dbAttemptsInt; attempts++ {
		log.Printf("Будет реализовано %d попыток доступа к файлу базы данных", dbAttemptsInt)
		exists, err = CheckIfDBExists(dbFile)
		// Проверяем
		if err == nil {
			log.Println("База данных найдена.")
			break // Успешная проверка, выходим из цикла
		} else {
			counter++
			log.Printf("Попытка %d: ошибка при проверке базы данных: %v", attempts+1, err)
			time.Sleep(1 * time.Second) // Ждем 1 секунду перед повторной попыткой
		}
	}
	// Если будет ошибка во всех попытках, то после завершения цикла переменная err будет содержать информацию об ошибке, возникшей на последней попытке.

	if err != nil {
		log.Printf("Не удалось проверить наличие базы данных после %d попыток: %v", counter, err)
		return dbFile
	}

	switch exists {
	case false:
		log.Println("База данных не существует")
		log.Println("Переходим к созданию директории для базы данных...")
		CreateDBDirectory(dbFile)  // Создаем директорию
		CreateDatabaseFile(dbFile) // Создаем файл базы данных
	case true:
		log.Println("База данных", dbFile, "уже существует.")
	}

	return dbFile // Возвращаем значение dbFil
}
