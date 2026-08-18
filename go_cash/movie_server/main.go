package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type Movie struct {
	ID     int      `json:"id"`
	Title  string   `json:"title"`
	Year   int      `json:"year"`
	Rating float64  `json:"rating"`
	Genre  string   `json:"genre"`
	Actors []string `json:"actors"`
}

var db *pgxpool.Pool

func main() {
	godotenv.Load()

	connString := os.Getenv("DATABASE_URL")
	if connString == "" {
		log.Fatal("DATABASE_URL")
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
	log.Println(".................Connected..................")

	app := fiber.New()
	app.Use("/", static.New("./public"))

	app.Get("/movies", getAllMovies)
	app.Get("/movies/:id", getMovieByID)
	app.Post("/movies", createMovie)
	app.Delete("/movies/:id", deleteMovie)
	app.Put("/movies/reorder", reorderMovies)
	app.Get("/filter", filterMovies)
	app.Get("/top", getTopMovies)

	app.Listen(":8000")
}

const selectCols = "SELECT id, title, year, rating, genre, actors FROM movies"
const defaultOrder = " ORDER BY position"

func errJSON(c fiber.Ctx, status int, msg string) error {
	return c.Status(status).JSON(fiber.Map{"error": msg})
}

func paramID(c fiber.Ctx) (id int, ok bool) {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		errJSON(c, 400, "invalid id")
		return 0, false
	}
	return id, true
}

func queryMovies(sql string, args ...any) ([]Movie, error) {
	rows, err := db.Query(context.Background(), sql, args...)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[Movie])
}

func firstOrErr(results []Movie, err error) (Movie, error) {
	if err == nil && len(results) == 0 {
		err = pgx.ErrNoRows
	}
	if len(results) > 0 {
		return results[0], nil
	}
	return Movie{}, err
}

func respondMovies(c fiber.Ctx, results []Movie, err error, notFoundMsg string) error {
	if err != nil {
		log.Println("query error:", err)
		return errJSON(c, 500, "database error")
	}
	if len(results) == 0 {
		return errJSON(c, 404, notFoundMsg)
	}
	return c.JSON(results)
}

func getAllMovies(c fiber.Ctx) error {
	results, err := queryMovies(selectCols + defaultOrder)
	return respondMovies(c, results, err, "no movies found")
}

func getMovieByID(c fiber.Ctx) error {
	id, ok := paramID(c)
	if !ok {
		return nil
	}
	movie, err := firstOrErr(queryMovies(selectCols+" WHERE id = $1", id))
	if err != nil {
		return errJSON(c, 404, "movie not found")
	}
	return c.JSON(movie)
}

func createMovie(c fiber.Ctx) error {
	var input Movie
	if err := c.Bind().Body(&input); err != nil {
		return errJSON(c, 400, "invalid request body")
	}
	if strings.TrimSpace(input.Title) == "" {
		return errJSON(c, 400, "title is required")
	}

	created, err := firstOrErr(queryMovies(`
		INSERT INTO movies (title, year, rating, genre, actors, position)
		VALUES ($1, $2, $3, $4, $5, COALESCE((SELECT MAX(position) FROM movies), 0) + 1)
		RETURNING id, title, year, rating, genre, actors`,
		input.Title, input.Year, input.Rating, input.Genre, input.Actors,
	))
	if err != nil {
		log.Println("create error:", err)
		return errJSON(c, 500, "could not create movie")
	}
	return c.Status(201).JSON(created)
}

func deleteMovie(c fiber.Ctx) error {
	id, ok := paramID(c)
	if !ok {
		return nil
	}

	tag, err := db.Exec(context.Background(), "DELETE FROM movies WHERE id = $1", id)
	if err != nil {
		log.Println("delete error:", err)
		return errJSON(c, 500, "could not delete movie")
	}
	if tag.RowsAffected() == 0 {
		return errJSON(c, 404, "movie not found")
	}
	return c.JSON(fiber.Map{"message": "movie deleted"})
}

func filterMovies(c fiber.Ctx) error {
	var clauses []string
	var args []any

	add := func(clause string, val any) {
		args = append(args, val)
		clauses = append(clauses, fmt.Sprintf(clause, len(args)))
	}

	if v := strings.TrimSpace(c.Query("title")); v != "" {
		add("title ILIKE $%d", "%"+v+"%")
	}
	if v := strings.TrimSpace(c.Query("actor")); v != "" {
		add("EXISTS (SELECT 1 FROM unnest(actors) AS a WHERE a ILIKE $%d)", "%"+v+"%")
	}
	if v := strings.TrimSpace(c.Query("genre")); v != "" {
		add("genre ILIKE $%d", v)
	}
	if v := strings.TrimSpace(c.Query("min")); v != "" {
		min, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return errJSON(c, 400, "invalid rating value")
		}
		add("rating >= $%d", min)
	}

	if len(clauses) == 0 {
		return errJSON(c, 400, "provide at least one filter: title, actor, genre, or min")
	}

	sql := selectCols + " WHERE " + strings.Join(clauses, " AND ") + " ORDER BY id"
	results, err := queryMovies(sql, args...)
	return respondMovies(c, results, err, "no movies found")
}

func getTopMovies(c fiber.Ctx) error {
	results, err := queryMovies(selectCols + " ORDER BY rating DESC LIMIT 10")
	return respondMovies(c, results, err, "no movies found")
}

func reorderMovies(c fiber.Ctx) error {
	var input struct {
		Order []int `json:"order"`
	}
	if err := c.Bind().Body(&input); err != nil || len(input.Order) == 0 {
		return errJSON(c, 400, "expected a non-empty \"order\" array of movie ids")
	}

	tx, err := db.Begin(context.Background())
	if err != nil {
		log.Println("reorder begin error:", err)
		return errJSON(c, 500, "could not reorder movies")
	}
	defer tx.Rollback(context.Background())

	for i, id := range input.Order {
		if _, err := tx.Exec(context.Background(),
			"UPDATE movies SET position = $1 WHERE id = $2", i, id); err != nil {
			log.Println("reorder update error:", err)
			return errJSON(c, 500, "could not reorder movies")
		}
	}

	if err := tx.Commit(context.Background()); err != nil {
		log.Println("reorder commit error:", err)
		return errJSON(c, 500, "could not reorder movies")
	}
	return c.JSON(fiber.Map{"message": "order updated"})
}
