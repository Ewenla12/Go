package main

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v3"
)

type Movie struct {
	ID     int      `json:"id"`
	Title  string   `json:"title"`
	Year   int      `json:"year"`
	Rating float64  `json:"rating"`
	Genre  string   `json:"genre"`
	Actors []string `json:"actors"`
}

// dummy data

var movies = []Movie{
	{1, "The Shawshank Redemption", 1994, 9.3, "Drama", []string{"Tim Robbins", "Morgan Freeman", "Bob Gunton"}},
	{2, "The Godfather", 1972, 9.2, "Crime", []string{"Marlon Brando", "Al Pacino", "James Caan"}},
	{3, "The Dark Knight", 2008, 9.0, "Action", []string{"Christian Bale", "Heath Ledger", "Aaron Eckhart"}},
	{4, "Pulp Fiction", 1994, 8.9, "Crime", []string{"John Travolta", "Uma Thurman", "Samuel L. Jackson"}},
	{5, "Forrest Gump", 1994, 8.8, "Drama", []string{"Tom Hanks", "Robin Wright", "Gary Sinise"}},
	{6, "Inception", 2010, 8.8, "Sci-Fi", []string{"Leonardo DiCaprio", "Joseph Gordon-Levitt", "Elliot Page"}},
	{7, "The Matrix", 1999, 8.7, "Sci-Fi", []string{"Keanu Reeves", "Laurence Fishburne", "Carrie-Anne Moss"}},
	{8, "Goodfellas", 1990, 8.7, "Crime", []string{"Robert De Niro", "Ray Liotta", "Joe Pesci"}},
	{9, "Fight Club", 1999, 8.8, "Drama", []string{"Brad Pitt", "Edward Norton", "Helena Bonham Carter"}},
	{10, "Interstellar", 2014, 8.7, "Sci-Fi", []string{"Matthew McConaughey", "Anne Hathaway", "Jessica Chastain"}},
	{11, "The Silence of the Lambs", 1991, 8.6, "Thriller", []string{"Jodie Foster", "Anthony Hopkins", "Scott Glenn"}},
	{12, "Schindler's List", 1993, 9.0, "History", []string{"Liam Neeson", "Ben Kingsley", "Ralph Fiennes"}},
	{13, "The Lord of the Rings", 2003, 9.0, "Fantasy", []string{"Elijah Wood", "Ian McKellen", "Viggo Mortensen"}},
	{14, "Star Wars: A New Hope", 1977, 8.6, "Sci-Fi", []string{"Mark Hamill", "Harrison Ford", "Carrie Fisher"}},
	{15, "Gladiator", 2000, 8.5, "Action", []string{"Russell Crowe", "Joaquin Phoenix", "Connie Nielsen"}},
	{16, "The Lion King", 1994, 8.5, "Animation", []string{"Matthew Broderick", "Jeremy Irons", "James Earl Jones"}},
	{17, "Avengers: Endgame", 2019, 8.4, "Action", []string{"Robert Downey Jr.", "Chris Evans", "Mark Ruffalo"}},
	{18, "Titanic", 1997, 7.9, "Romance", []string{"Leonardo DiCaprio", "Kate Winslet", "Billy Zane"}},
	{19, "Joker", 2019, 8.4, "Drama", []string{"Joaquin Phoenix", "Robert De Niro", "Zazie Beetz"}},
	{20, "Parasite", 2019, 8.5, "Thriller", []string{"Song Kang-ho", "Lee Sun-kyun", "Cho Yeo-jeong"}},
	{21, "Whiplash", 2014, 8.5, "Drama", []string{"Miles Teller", "J.K. Simmons", "Melissa Benoist"}},
	{22, "La La Land", 2016, 8.0, "Romance", []string{"Ryan Gosling", "Emma Stone", "John Legend"}},
	{23, "Get Out", 2017, 7.7, "Horror", []string{"Daniel Kaluuya", "Allison Williams", "Bradley Whitford"}},
	{24, "Black Panther", 2018, 7.3, "Action", []string{"Chadwick Boseman", "Michael B. Jordan", "Lupita Nyong'o"}},
	{25, "Spider-Man: Into the Spider-Verse", 2018, 8.4, "Animation", []string{"Shameik Moore", "Jake Johnson", "Hailee Steinfeld"}},
	{26, "1917", 2019, 8.3, "War", []string{"George MacKay", "Dean-Charles Chapman", "Mark Strong"}},
	{27, "Baby Driver", 2017, 7.6, "Action", []string{"Ansel Elgort", "Jon Hamm", "Jamie Foxx"}},
	{28, "Mad Max: Fury Road", 2015, 8.1, "Action", []string{"Tom Hardy", "Charlize Theron", "Nicholas Hoult"}},
	{29, "The Revenant", 2015, 8.0, "Adventure", []string{"Leonardo DiCaprio", "Tom Hardy", "Will Poulter"}},
	{30, "Gone Girl", 2014, 8.1, "Thriller", []string{"Ben Affleck", "Rosamund Pike", "Neil Patrick Harris"}},
	{31, "Arrival", 2016, 7.9, "Sci-Fi", []string{"Amy Adams", "Jeremy Renner", "Forest Whitaker"}},
	{32, "Hereditary", 2018, 7.3, "Horror", []string{"Toni Collette", "Milly Shapiro", "Gabriel Byrne"}},
	{33, "Midsommar", 2019, 7.1, "Horror", []string{"Florence Pugh", "Jack Reynor", "William Jackson Harper"}},
	{34, "The Irishman", 2019, 7.8, "Crime", []string{"Robert De Niro", "Al Pacino", "Joe Pesci"}},
	{35, "Once Upon a Time in Hollywood", 2019, 7.6, "Drama", []string{"Leonardo DiCaprio", "Brad Pitt", "Margot Robbie"}},
	{36, "Knives Out", 2019, 7.9, "Mystery", []string{"Daniel Craig", "Chris Evans", "Ana de Armas"}},
	{37, "The Grand Budapest Hotel", 2014, 8.1, "Comedy", []string{"Ralph Fiennes", "Tony Revolori", "Saoirse Ronan"}},
	{38, "Dune", 2021, 8.0, "Sci-Fi", []string{"Timothée Chalamet", "Rebecca Ferguson", "Oscar Isaac"}},
	{39, "No Time to Die", 2021, 7.3, "Action", []string{"Daniel Craig", "Rami Malek", "Léa Seydoux"}},
	{40, "The Batman", 2022, 7.8, "Action", []string{"Robert Pattinson", "Zoë Kravitz", "Paul Dano"}},
	{41, "Everything Everywhere All at Once", 2022, 7.8, "Sci-Fi", []string{"Michelle Yeoh", "Ke Huy Quan", "Jamie Lee Curtis"}},
	{42, "Top Gun: Maverick", 2022, 8.3, "Action", []string{"Tom Cruise", "Miles Teller", "Jennifer Connelly"}},
	{43, "The Northman", 2022, 7.1, "Adventure", []string{"Alexander Skarsgård", "Nicole Kidman", "Anya Taylor-Joy"}},
	{44, "Nope", 2022, 6.9, "Horror", []string{"Daniel Kaluuya", "Keke Palmer", "Steven Yeun"}},
	{45, "Avatar", 2009, 7.9, "Sci-Fi", []string{"Sam Worthington", "Zoe Saldana", "Sigourney Weaver"}},
	{46, "Oppenheimer", 2023, 8.9, "History", []string{"Cillian Murphy", "Emily Blunt", "Matt Damon"}},
	{47, "Barbie", 2023, 7.0, "Comedy", []string{"Margot Robbie", "Ryan Gosling", "America Ferrera"}},
	{48, "John Wick", 2014, 7.4, "Action", []string{"Keanu Reeves", "Michael Nyqvist", "Alfie Allen"}},
	{49, "The Prestige", 2006, 8.5, "Mystery", []string{"Christian Bale", "Hugh Jackman", "Scarlett Johansson"}},
	{50, "Memento", 2000, 8.4, "Mystery", []string{"Guy Pearce", "Carrie-Anne Moss", "Joe Pantoliano"}},
}

func main() {
	app := fiber.New()

	// all movies
	app.Get("/movies", getAllMovies)

	// movie by id
	app.Get("/movies/:id", getMovieByID)

	// search by title
	app.Get("/search", searchByTitle)

	// search by actor
	app.Get("/actors", searchByActor)

	// filter by genre
	app.Get("/genre", filterByGenre)

	// filter by rating
	app.Get("/rating", filterByRating)

	// top 10 highest rated
	app.Get("/top", getTopMovies)

	app.Listen(":8000")
}

//  functions

// show all movies
func getAllMovies(c fiber.Ctx) error {
	return c.JSON(movies)
}

// show one movie by id
func getMovieByID(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid id",
		})
	}
	for _, m := range movies {
		if m.ID == id {
			return c.JSON(m)
		}
	}
	return c.Status(404).JSON(fiber.Map{
		"error": "movie not found",
	})
}

// search movies by title
func searchByTitle(c fiber.Ctx) error {
	query := strings.ToLower(c.Query("title"))
	if query == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "provide a title query eg. /search?title=dark",
		})
	}
	var results []Movie
	for _, m := range movies {
		if strings.Contains(strings.ToLower(m.Title), query) {
			results = append(results, m)
		}
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
	query := strings.ToLower(c.Query("name"))
	if query == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "provide a name query eg. /actors?name=morgan",
		})
	}
	var results []Movie
	for _, m := range movies {
		for _, actor := range m.Actors {
			if strings.Contains(strings.ToLower(actor), query) {
				results = append(results, m)
				break
			}
		}
	}
	if len(results) == 0 {
		return c.Status(404).JSON(fiber.Map{
			"error": "no movies found for that actor",
		})
	}
	return c.JSON(results)
}

// filter by genre
func filterByGenre(c fiber.Ctx) error {
	query := strings.ToLower(c.Query("name"))
	if query == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "provide a genre eg. /genre?name=action",
		})
	}
	var results []Movie
	for _, m := range movies {
		if strings.ToLower(m.Genre) == query {
			results = append(results, m)
		}
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
	var results []Movie
	for _, m := range movies {
		if m.Rating >= min {
			results = append(results, m)
		}
	}
	if len(results) == 0 {
		return c.Status(404).JSON(fiber.Map{
			"error": "no movies found with that rating",
		})
	}
	return c.JSON(results)
}

// show top 10 highest rated movies
func getTopMovies(c fiber.Ctx) error {
	sorted := make([]Movie, len(movies))
	copy(sorted, movies)

	//used bubble sort from c++
	for i := 0; i < len(sorted); i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[j].Rating > sorted[i].Rating {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}

	top := 10
	if len(sorted) < top {
		top = len(sorted)
	}
	return c.JSON(sorted[:top])
}
