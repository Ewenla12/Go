package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Movie struct {
	ID     int      `json:"id"`
	Title  string   `json:"title"`
	Year   int      `json:"year"`
	Rating float64  `json:"rating"`
	Genre  string   `json:"genre"`
	Actors []string `json:"actors"`
}

// db is a package-level Postgres connection pool, shared by every handler.
var db *pgxpool.Pool

func main() {
	connString := os.Getenv("DATABASE_URL")
	if connString == "" {
		connString = "postgres://postgres:154993@localhost:5432/movieapp"
	}

	var err error
	db, err = pgxpool.New(context.Background(), connString)
	if err != nil {
		log.Fatal("unable to connect to database: ", err)
	}
	defer db.Close()

	if err := db.Ping(context.Background()); err != nil {
		log.Fatal("database ping failed: ", err)
	}
	log.Println("connected to postgres")
	app := fiber.New()

	// connect frontend
	app.Use("/", static.New("./public"))

	// all movies
	app.Get("/movies", getAllMovies)

	// search by id
	app.Get("/movies/:id", getMovieByID)

	// add a new movie
	app.Post("/movies", createMovie)

	// remove a movie
	app.Delete("/movies/:id", deleteMovie)

	// search by title
	app.Get("/search", searchByTitle)

	// search by actor
	app.Get("/actors", searchByActor)

	// search by genre
	app.Get("/genre", filterByGenre)

	// filter by rating
	app.Get("/rating", filterByRating)

	// top 10 rated
	app.Get("/top", getTopMovies)

	app.Listen(":8000")
}

// functions
func scanMovies(rows interface {
	Next() bool
	Scan(dest ...any) error
	Err() error
}) ([]Movie, error) {
	var results []Movie
	for rows.Next() {
		var m Movie
		if err := rows.Scan(&m.ID, &m.Title, &m.Year, &m.Rating, &m.Genre, &m.Actors); err != nil {
			return nil, err
		}
		results = append(results, m)
	}
	return results, rows.Err()
}

// show all movies
func getAllMovies(c fiber.Ctx) error {
	rows, err := db.Query(context.Background(),
		"SELECT id, title, year, rating, genre, actors FROM movies ORDER BY id")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "database error"})
	}
	defer rows.Close()

	results, err := scanMovies(rows)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "database error"})
	}
	return c.JSON(results)
}

// search movie by id
func getMovieByID(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid id",
		})
	}

	var m Movie
	err = db.QueryRow(context.Background(),
		"SELECT id, title, year, rating, genre, actors FROM movies WHERE id = $1", id,
	).Scan(&m.ID, &m.Title, &m.Year, &m.Rating, &m.Genre, &m.Actors)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "movie not found",
		})
	}
	return c.JSON(m)
}

func createMovie(c fiber.Ctx) error {
	var input Movie
	if err := c.Bind().Body(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	if strings.TrimSpace(input.Title) == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "title is required",
		})
	}

	var created Movie
	err := db.QueryRow(context.Background(), `
		INSERT INTO movies (title, year, rating, genre, actors)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, title, year, rating, genre, actors`,
		input.Title, input.Year, input.Rating, input.Genre, input.Actors,
	).Scan(&created.ID, &created.Title, &created.Year, &created.Rating, &created.Genre, &created.Actors)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "could not create movie"})
	}

	return c.Status(201).JSON(created)
}

// deleteMovie removes a movie by id.
func deleteMovie(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid id",
		})
	}

	tag, err := db.Exec(context.Background(), "DELETE FROM movies WHERE id = $1", id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "could not delete movie"})
	}
	if tag.RowsAffected() == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "movie not found"})
	}

	return c.Status(200).JSON(fiber.Map{"message": "movie deleted"})
}

// search movie by title
func searchByTitle(c fiber.Ctx) error {
	query := strings.TrimSpace(c.Query("title"))
	if query == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "provide a title query eg. /search?title=dark",
		})
	}

	rows, err := db.Query(context.Background(),
		"SELECT id, title, year, rating, genre, actors FROM movies WHERE title ILIKE $1 ORDER BY id",
		"%"+query+"%",
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "database error"})
	}
	defer rows.Close()

	results, err := scanMovies(rows)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "database error"})
	}
	if len(results) == 0 {
		return c.Status(404).JSON(fiber.Map{
			"error": "no movies found",
		})
	}
	return c.JSON(results)
}

// search movies by actor name
func searchByActor(c fiber.Ctx) error {
	query := strings.TrimSpace(c.Query("name"))
	if query == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "provide a name query eg. /actors?name=morgan",
		})
	}

	rows, err := db.Query(context.Background(), `
		SELECT id, title, year, rating, genre, actors FROM movies m
		WHERE EXISTS (
			SELECT 1 FROM unnest(m.actors) AS actor
			WHERE actor ILIKE $1
		)
		ORDER BY id`,
		"%"+query+"%",
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "database error"})
	}
	defer rows.Close()

	results, err := scanMovies(rows)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "database error"})
	}
	if len(results) == 0 {
		return c.Status(404).JSON(fiber.Map{
			"error": "no movies found for that actor",
		})
	}
	return c.JSON(results)
}

// search by genre
func filterByGenre(c fiber.Ctx) error {
	query := strings.TrimSpace(c.Query("name"))
	if query == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "provide a genre eg. /genre?name=action",
		})
	}

	rows, err := db.Query(context.Background(),
		"SELECT id, title, year, rating, genre, actors FROM movies WHERE genre ILIKE $1 ORDER BY id",
		query,
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "database error"})
	}
	defer rows.Close()

	results, err := scanMovies(rows)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "database error"})
	}
	if len(results) == 0 {
		return c.Status(404).JSON(fiber.Map{
			"error": "no movies found for that genre",
		})
	}
	return c.JSON(results)
}

// to get rating
func filterByRating(c fiber.Ctx) error {
	query := c.Query("min")
	if query == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "provide a min rating eg. /rating?min=8.5",
		})
	}
	min, err := strconv.ParseFloat(query, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid rating value",
		})
	}

	rows, err := db.Query(context.Background(),
		"SELECT id, title, year, rating, genre, actors FROM movies WHERE rating >= $1 ORDER BY id",
		min,
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "database error"})
	}
	defer rows.Close()

	results, err := scanMovies(rows)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "database error"})
	}
	if len(results) == 0 {
		return c.Status(404).JSON(fiber.Map{
			"error": "no movies found with that rating",
		})
	}
	return c.JSON(results)
}

// show top 10 rated movies
func getTopMovies(c fiber.Ctx) error {
	rows, err := db.Query(context.Background(),
		"SELECT id, title, year, rating, genre, actors FROM movies ORDER BY rating DESC LIMIT 10")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "database error"})
	}
	defer rows.Close()

	results, err := scanMovies(rows)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "database error"})
	}
	return c.JSON(results)
}
