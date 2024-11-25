package database

import (
	"database/sql"
	"log"
)

// openDatabase открывает соединение с базой данных.
// С помощью sql.Open открывается база данных SQLite.
// Если база данных не существует, то файл будет создан.
func OpenDatabase(dbFile string) *sql.DB {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatalf("Фатальная ошибка открытия базы данных: %v\n", err)
	}
	log.Println("Соединение с базой данных успешно установлено:", dbFile)
	return db
}

// createTable создает таблицу и индекс в базе данных.
func CreateTable(db *sql.DB) {
	createTableSQL := `
    CREATE TABLE scheduler (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        date TEXT NOT NULL,
        title TEXT NOT NULL,
        comment TEXT,
        repeat TEXT CHECK(length(repeat) <= 128)
    );`

	createIndexSQL := `CREATE INDEX idx_date ON scheduler (date);`

	// Выполняем SQL-запрос для создания таблицы
	if _, err := db.Exec(createTableSQL); err != nil {
		log.Fatalf("Фатальная ошибка при создании таблицы: %v", err)
	}
	log.Println("Таблица scheduler успешно создана.")

	// Выполняем SQL-запрос для создания индекса
	if _, err := db.Exec(createIndexSQL); err != nil {
		log.Fatalf("Фатальная ошибка при создании индекса: %v", err)
	}
	log.Println("Индекс idx_date успешно создан.")
}
