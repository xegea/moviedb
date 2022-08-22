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
	router.Handle("/img/favicon.png", fs)
	router.Handle("/img/logo.png", fs)
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
