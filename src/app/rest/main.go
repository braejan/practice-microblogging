package main

import (
	"context"
	"log"
	"net/http"
	"os"

	microblogHandler "github.com/braejan/practice-microblogging/src/app/rest/handlers/microblog"
	userHandler "github.com/braejan/practice-microblogging/src/app/rest/handlers/user"
	"github.com/braejan/practice-microblogging/src/app/rest/server"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	printErr(err, true)
	configuration := &server.Config{
		Port: os.Getenv("PORT"),
	}
	server, err := server.NewServer(context.Background(), configuration)
	printErr(err, true)
	server.Start(BindRoutes)

}

func BindRoutes(server server.Server, r *mux.Router) {
	r.Handle("/users", userHandler.CreateUserHandler(server)).Methods(http.MethodPost)
	r.Handle("/users/{userID}", userHandler.GetUserHandler(server)).Methods(http.MethodGet)
	r.Handle("/microblog/all", microblogHandler.GetAllPostHandler(server)).Methods(http.MethodGet)
	r.Handle("/microblog", microblogHandler.CreatePostHandler(server)).Methods(http.MethodPost)
	r.Handle("/microblog/{postID}", microblogHandler.GetPostByIDHandler(server)).Methods(http.MethodGet)
	r.Handle("/microblog/user/{userID}", microblogHandler.GetAllPostByUserIDHandler(server)).Methods(http.MethodGet)
	r.Handle("/microblog/{postID}/like", microblogHandler.LikePostHandler(server)).Methods(http.MethodPut)
	r.Handle("/microblog/{postID}/dislike", microblogHandler.DislikePostHandler(server)).Methods(http.MethodPut)
}

func printErr(err error, exit bool) {
	if err != nil {
		if exit {
			log.Fatalf("fatal error: %v", err)
		}
		log.Printf("error: %v", err)
	}
}
