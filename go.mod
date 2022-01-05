module github.com/moviedb/api

// +heroku goVersion go1.17
go 1.17

require github.com/gorilla/mux v1.8.0

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/felixge/httpsnoop v1.0.1 // indirect
)

require (
	github.com/RediSearch/redisearch-go v1.1.1
	github.com/gomodule/redigo v1.8.3
	github.com/gorilla/handlers v1.5.1
	github.com/joho/godotenv v1.4.0
)
