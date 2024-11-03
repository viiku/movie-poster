// The JSON-based web service of the Open Movie Database lets you search
// https://omdbapi.com/ for a movie by name and download its poster image.

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type Rating struct {
	Source string `json:"Source"`
	Value  string `json:"Value"`
}

type Movie struct {
	Title    string `json:"Title"`
	Year     string `json:"Year"`
	Rated    string `json:"Rated"`
	Released string `json:"Released"`
	Director string `json:"Director"`
	Language string `json:"Language"`
	Country  string `json:"Country"`
	Poster   string `json:"Poster"`
	Ratings  []Rating
	Error    string `json:"Error"`
}

func getAPIKey() string {
	apikey := os.Getenv("OMDB_API_KEY")
	if apikey == "" {
		log.Fatal("OMDB api key not set")
	}
	return apikey
}

func fetchMovieData(movieName string) (Movie, error) {
	apikey := getAPIKey()
	url := fmt.Sprintf("https://www.omdbapi.com/?apikey=%s&t=%s", apikey, movieName)

	resp, err := http.Get(url)
	if err != nil {
		return Movie{}, fmt.Errorf("failed to make request: %w", err)
	}

	var movieData Movie
	if err := json.NewDecoder(resp.Body).Decode(&movieData); err != nil {
		return Movie{}, fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Check if OMDb returned an error message
	if movieData.Error != "" {
		return movieData, fmt.Errorf("OMDb API error: %s", movieData.Error)
	}

	return movieData, nil
}

func downloadMoviePoster(posterUrl string, fileName string) error {
	resp, err := http.Get(posterUrl)
	if err != nil {
		return fmt.Errorf("failed to download poster: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return fmt.Errorf("downloading poster failed: %s", resp.Status)
	}

	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return fmt.Errorf("creating file failed: %s", resp.Status)
	}

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to save poster: %w", err)
	}

	return nil
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go <movie name>")
	}
	movieName := strings.Join(os.Args[1:], " ")

	// Fetch movie data from OMDb
	movie, err := fetchMovieData(movieName)
	if err != nil {
		log.Fatalf("Error fetching movie data: %v", err)
	}
	fmt.Print(movie)

	// Download the poster
	// if movie.Poster == "N/A" {
	// 	log.Fatalf("No poster available for %s", movie.Title)
	// }
	// filename := movie.Title + ".jpg"
	// err = downloadMoviePoster(movie.Poster, filename)
	// if err != nil {
	// 	log.Fatalf("Error downloading poster: %v", err)
	// }

	// fmt.Printf("Poster for %s downloaded as %s\n", movie.Title, filename)
}
