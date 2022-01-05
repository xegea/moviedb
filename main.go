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

type Content struct {
	Title         string    `json:",omitempty"`
	Url           string    `json:",omitempty"`
	ContentRating string    `json:",omitempty"`
	Type          string    `json:",omitempty"`
	Description   string    `json:",omitempty"`
	Genre         string    `json:",omitempty"`
	Image         string    `json:",omitempty"`
	DateCreated   int64     `json:",omitempty"`
	Director      []string  `json:",omitempty"`
	Actors        []string  `json:",omitempty"`
	Trailer       []Trailer `json:",omitempty"`
}

type Trailer struct {
	Name         string `json:",omitempty"`
	Description  string `json:",omitempty"`
	Url          string `json:",omitempty"`
	ThumbnailUrl string `json:",omitempty"`
}

var (
	client *redisearch.Client

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
	client = redisearch.NewClientFromPool(pool, "idx:title")

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

	docs, total, err := client.Search(redisearch.NewQuery(fmt.Sprint("@title:", query, "*")).
		// SetReturnFields("title", "description", "type"). // if SetReturnFields
		Limit(page*10, 10))
	if err != nil {
		log.Fatal(err)
	}

	var contentList []Content

	for _, v := range docs {

		var jsonbody string
		for _, p := range v.Properties {
			jsonbody = fmt.Sprint(p)
		}
		if err != nil {
			log.Fatal(err)
		}

		var content Content
		if err := json.Unmarshal([]byte(jsonbody), &content); err != nil {
			log.Fatalf("Failed to Unmarshall %s", jsonbody)
		}
		contentList = append(contentList, content)
	}

	json.NewEncoder(w).Encode(contentList)

	fmt.Printf("total: %d\n", total)
}
