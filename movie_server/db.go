package main

import (
	"database/sql"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/lib/pq"
)

var db *sql.DB

func initDB() {
	var err error
	// adjust user/password/dbname to match your setup
	connStr := "postgres://postgres:154993localhost:5432/moviesdb?sslmode=disable"
	db, err = sql.Open("pgx", connStr)
	if err != nil {
		log.Fatal("failed to open db:", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal("failed to connect to db:", err)
	}

	createTable := `
	CREATE TABLE IF NOT EXISTS movies (
		id SERIAL PRIMARY KEY,
		title TEXT NOT NULL,
		year INT NOT NULL,
		rating FLOAT8 NOT NULL,
		genre TEXT NOT NULL,
		actors TEXT[] NOT NULL
	);`
	if _, err = db.Exec(createTable); err != nil {
		log.Fatal("failed to create table:", err)
	}

	seedIfEmpty()
}

func seedIfEmpty() {
	var count int
	db.QueryRow("SELECT COUNT(*) FROM movies").Scan(&count)
	if count > 0 {
		return
	}
	log.Println("seeding movies table...")
	for _, m := range movies {
		_, err := db.Exec(
			`INSERT INTO movies (title, year, rating, genre, actors) VALUES ($1, $2, $3, $4, $5)`,
			m.Title, m.Year, m.Rating, m.Genre, pq.Array(m.Actors),
		)
		if err != nil {
			log.Println("seed error:", err)
		}
	}
}
