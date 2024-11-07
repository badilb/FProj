package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"net/http"
	"strconv"

	_ "github.com/lib/pq" // Импортируем драйвер PostgreSQL
)

// Структура записи
type Record struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}

// Подключение к базе данных
const (
	host     = "localhost"
	port     = 5482
	user     = "postgres"
	password = "postgres"   // Укажите свой пароль
	dbname   = "ast_census" // Укажите имя вашей базы данных
)

var db *sql.DB

// Инициализация базы данных
func initDB() {
	var err error
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		panic(err)
	}

	fmt.Println("Подключено к базе данных!")
}

// Получить все записи
func getRecords(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	rows, err := db.Query("SELECT id, text FROM records")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var records []Record
	for rows.Next() {
		var record Record
		if err := rows.Scan(&record.ID, &record.Text); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		records = append(records, record)
	}

	json.NewEncoder(w).Encode(records)
}

// Добавить запись
func addRecord(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var record Record
	if err := json.NewDecoder(r.Body).Decode(&record); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := db.QueryRow("INSERT INTO records(text) VALUES($1) RETURNING id", record.Text).Scan(&record.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(record)
}

// Обновить запись
func updateRecord(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	var record Record
	if err := json.NewDecoder(r.Body).Decode(&record); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = db.Exec("UPDATE records SET text=$1 WHERE id=$2", record.Text, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	record.ID = id
	json.NewEncoder(w).Encode(record)
}

// Удалить запись
func deleteRecord(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	_, err = db.Exec("DELETE FROM records WHERE id=$1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent) // Успешно, без содержимого
}

func main() {
	initDB()
	defer db.Close()

	router := mux.NewRouter()
	router.HandleFunc("/records", getRecords).Methods("GET")
	router.HandleFunc("/records", addRecord).Methods("POST")
	router.HandleFunc("/records/{id}", updateRecord).Methods("PUT")
	router.HandleFunc("/records/{id}", deleteRecord).Methods("DELETE")

	// CORS middleware
	handler := cors.Default().Handler(router)

	http.Handle("/", handler)

	// Сообщение о запуске сервера
	fmt.Println("Сервер запущен на http://localhost:9070")
	if err := http.ListenAndServe(":9070", nil); err != nil {
		fmt.Println("Ошибка запуска сервера:", err)
	}
}
