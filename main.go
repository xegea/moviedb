package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type Movie struct {
	Title         map[string]string `json:",omitempty"`
	Url           string            `json:",omitempty"`
	ContentRating string            `json:",omitempty"`
	Type          string            `json:",omitempty"`
	Description   map[string]string `json:",omitempty"`
	Genre         string            `json:",omitempty"`
	Image         string            `json:",omitempty"`
	DateCreated   int64             `json:",omitempty"`
	Director      []string          `json:",omitempty"`
	Actors        []string          `json:",omitempty"`
	Trailer       []Trailer         `json:",omitempty"`
}

type Trailer struct {
	Name         map[string]string `json:",omitempty"`
	Description  map[string]string `json:",omitempty"`
	Url          string            `json:",omitempty"`
	ThumbnailUrl string            `json:",omitempty"`
}

var (
	client_es *redisearch.Client
	client_us *redisearch.Client

	//go:embed static/*
	static embed.FS
)

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("Error loading .env file")
	}

	host := os.Getenv("REDIS_HOST")
	password := os.Getenv("REDIS_PASSWORD")

	pool := &redis.Pool{Dial: func() (redis.Conn, error) {
		return redis.Dial("tcp", host, redis.DialPassword(password))
	}}
	client_es = redisearch.NewClientFromPool(pool, "idx:title:es")
	client_us = redisearch.NewClientFromPool(pool, "idx:title:us")

	handleRequests()
}

func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)

	fsys, err := fs.Sub(static, "static")
	if err != nil {
		panic(err)
	}

	// Serve static files
	fs := http.FileServer(http.FS(fsys))
	router.Handle("/", fs)
	router.Handle("/js/index.js", fs)
	router.Handle("/css/style.css", fs)

	router.HandleFunc("/search/", searchHandler)

	port := ":" + os.Getenv("PORT")
	if os.Getenv("PORT") == "" {
		port = ":8080"
	}

	cors := handlers.AllowedOrigins([]string{"*"})
	fmt.Println("Listen on port", port)
	log.Fatal(http.ListenAndServe(port, handlers.CORS(cors)(router)))
}

func searchHandler(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query().Get("q")
	var page int
	fmt.Sscan(r.URL.Query().Get("p"), &page)
	country := r.URL.Query().Get("c")

	client := resolveClient(country)

	docs, total, err := client.Search(redisearch.NewQuery(fmt.Sprint("@title:", query, "*")).
		// SetReturnFields("title", "description", "type"). // if SetReturnFields
		Limit(page*10, 10))
	if err != nil {
		log.Fatal(err)
	}

	var movieList []Movie

	for _, v := range docs {

		var jsonbody string
		for _, p := range v.Properties {
			jsonbody = fmt.Sprint(p)
		}
		if err != nil {
			log.Fatal(err)
		}

		var content Movie
		if err := json.Unmarshal([]byte(jsonbody), &content); err != nil {
			log.Fatalf("Failed to Unmarshall %s", jsonbody)
		}
		movieList = append(movieList, content)
	}

	json.NewEncoder(w).Encode(movieList)

	fmt.Printf("total: %d\n", total)
}

func resolveClient(country string) *redisearch.Client {
	var client *redisearch.Client
	switch country {
	case "es":
		{
			client = client_es
		}
	case "us":
		{
			client = client_us
		}
	default:
		{
			client = client_es
		}
	}
	return client
}
