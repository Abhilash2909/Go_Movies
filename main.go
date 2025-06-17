package main

// Importing necessary standard and third-party packages
import (
	"encoding/json" // To handle JSON encoding and decoding
	"fmt"
	"log"
	"net/http" // For HTTP server and request/response handling

	"github.com/google/uuid" // To generate unique IDs for each movie
	"github.com/gorilla/mux" // To use the Mux router for handling routes with parameters
)

// Movie struct defines the structure of a movie object
type Movie struct {
	ID       string    `json:"id"`       // Unique identifier for the movie
	Isbn     string    `json:"isbn"`     // ISBN number of the movie
	Title    string    `json:"title"`    // Title of the movie
	Director *Director `json:"director"` // Director details embedded as a pointer
}

// Director struct defines the structure for director details
type Director struct {
	Firstname string `json:"firstname"` // First name of the director
	Lastname  string `json:"lastname"`  // Last name of the director
}

// A global slice to act as an in-memory database for storing movies
var movies []Movie

// Handler to return all movies as JSON
func getMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") // Set response header to JSON
	json.NewEncoder(w).Encode(movies)                  // Encode and send all movies
}

// Handler to return a single movie by ID
func getMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") // Set response header to JSON
	params := mux.Vars(r)                              // Get route parameters from request (e.g., ID)
	for _, item := range movies {                      // Iterate over all movies
		if item.ID == params["id"] { // If movie ID matches
			json.NewEncoder(w).Encode(item) // Encode and return the movie
			return
		}
	}
	http.Error(w, "Movie not found", http.StatusNotFound) // Send 404 if not found
}

// Handler to delete a movie by ID
func deleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") // Set response header to JSON
	params := mux.Vars(r)                              // Get the ID from route parameters
	for index, item := range movies {                  // Iterate over movies
		if item.ID == params["id"] { // If movie ID matches
			// Remove movie by slicing around the index
			movies = append(movies[:index], movies[index+1:]...)
			json.NewEncoder(w).Encode(movies) // Return the updated movie list
			return
		}
	}
	http.Error(w, "Movie not found", http.StatusNotFound) // Send 404 if not found
}

// Handler to create a new movie entry
func createMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")             // Set response header to JSON
	var movie Movie                                                // Define a variable to hold incoming JSON
	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil { // Decode the request body
		http.Error(w, "Invalid JSON", http.StatusBadRequest) // Return 400 if bad JSON
		return
	}
	movie.ID = uuid.New().String()   // Generate a new unique ID for the movie
	movies = append(movies, movie)   // Add the new movie to the slice
	json.NewEncoder(w).Encode(movie) // Return the newly created movie
}

// Handler to update an existing movie by ID
func updateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") // Set response header to JSON
	params := mux.Vars(r)                              // Get the movie ID from route

	for index, item := range movies { // Iterate through the movies
		if item.ID == params["id"] { // Find the movie to update
			movies = append(movies[:index], movies[index+1:]...) // Delete old movie

			var movie Movie
			if err := json.NewDecoder(r.Body).Decode(&movie); err != nil { // Decode new data
				http.Error(w, "Invalid request body", http.StatusBadRequest) // Return error
				return
			}

			movie.ID = params["id"]          // Keep the same ID
			movies = append(movies, movie)   // Add the updated movie
			json.NewEncoder(w).Encode(movie) // Return the updated movie
			return
		}
	}
	http.Error(w, "Movie not found", http.StatusNotFound) // Send 404 if not found
}

// main function sets up the router and starts the server
func main() {
	r := mux.NewRouter() // Create a new Gorilla Mux router

	// Add some sample movies to the in-memory slice
	movies = append(movies, Movie{
		ID:    "1",
		Isbn:  "001",
		Title: "Movie One",
		Director: &Director{
			Firstname: "John",
			Lastname:  "Doe",
		},
	})
	movies = append(movies, Movie{
		ID:    "2",
		Isbn:  "002",
		Title: "Movie Two",
		Director: &Director{
			Firstname: "Jane",
			Lastname:  "Doe",
		},
	})

	// Define the API routes and map them to their respective handler functions
	r.HandleFunc("/movies", getMovies).Methods("GET")           // GET all movies
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")       // GET a single movie by ID
	r.HandleFunc("/movies", createMovie).Methods("POST")        // POST a new movie
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")    // PUT to update a movie
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE") // DELETE a movie

	fmt.Println("Starting server on port 8000...") // Print server start message
	log.Fatal(http.ListenAndServe(":8000", r))     // Start the server and listen on port 8000
}
