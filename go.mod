module github.com/moviedb/api

// +heroku goVersion go1.17
go 1.17

require github.com/gorilla/mux v1.8.0

require github.com/felixge/httpsnoop v1.0.1 // indirect

require (
	github.com/gorilla/handlers v1.5.1
	github.com/joho/godotenv v1.4.0
)
