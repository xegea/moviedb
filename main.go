package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/moviedb/api/pkg/config"
	"github.com/moviedb/api/pkg/handler"
	"github.com/moviedb/api/pkg/server"
)

type Movie struct {
	Title         map[string]string `json:",omitempty"`
	Url           string            `json:",omitempty"`
	ContentRating string            `json:",omitempty"`
	Type          string            `json:",omitempty"`
	Description   map[string]string `json:",omitempty"`
	Genre         string            `json:",omitempty"`
	Image         string            `json:",omitempty"`
	ReleaseDate   int64             `json:",omitempty"`
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
	//go:embed static/*
	static embed.FS
)

func main() {

	env := flag.String("env", ".env", "environment path")
	flag.Parse()

	cfg, err := config.LoadConfig(env)
	if err != nil {
		log.Fatalf("unable to load config: %+v", err)
	}

	router := mux.NewRouter().StrictSlash(true)

	srv := server.NewServer(
		cfg,
		router,
	)

	fsys, err := fs.Sub(static, "static")
	if err != nil {
		panic(err)
	}

	// Serve static files
	fs := http.FileServer(http.FS(fsys))
	router.Handle("/", fs)
	router.Handle("/js/index.js", fs)
	router.Handle("/css/style.css", fs)

	router.HandleFunc("/search/", handler.SearchHandler(srv))

	port := ":" + os.Getenv("PORT")
	if os.Getenv("PORT") == "" {
		port = ":8080"
	}

	cors := handlers.AllowedOrigins([]string{"*"})
	fmt.Println("Listen on port", port)
	log.Fatal(http.ListenAndServe(port, handlers.CORS(cors)(router)))
}
