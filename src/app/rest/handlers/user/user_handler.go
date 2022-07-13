package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/braejan/practice-microblogging/src/app/rest/server"
	"github.com/braejan/practice-microblogging/src/domain/user/entities"
	"github.com/braejan/practice-microblogging/src/domain/user/usecases"
	"github.com/gorilla/mux"
)

func CreateUserHandler(server server.Server) (handlerFunc http.HandlerFunc) {
	handlerFunc = func(w http.ResponseWriter, r *http.Request) {
		var requestUser entities.User
		err := json.NewDecoder(r.Body).Decode(&requestUser)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Printf("error getting user from body")
			return
		}
		userUsesCases, err := usecases.NewUserUsecases()
		if err != nil {
			log.Printf("error creating userUsecases: %s", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = userUsesCases.CreateUser(requestUser.ID, requestUser.Name)
		if err != nil {
			log.Printf("error creating user: %s", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
	}
	return
}

func GetUserHandler(server server.Server) (handlerFunc http.HandlerFunc) {
	handlerFunc = func(w http.ResponseWriter, r *http.Request) {
		userID := mux.Vars(r)["userID"]
		if userID == "" {
			w.WriteHeader(http.StatusBadRequest)
			log.Printf("error getting userID from path")
			return
		}
		requestUserID := r.Header.Get("user")
		if requestUserID == "" {
			w.WriteHeader(http.StatusBadRequest)
			log.Printf("error getting userID from header")
			return
		}
		userUsesCases, err := usecases.NewUserUsecases()
		if err != nil {
			log.Printf("error creating userUsecases: %s", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		user, err := userUsesCases.FindUserByID(userID)
		if err != nil {
			log.Printf("error creating user: %s", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if user == nil {
			log.Printf("user %s not found", userID)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(user)
		if err != nil {
			log.Printf("error getting userID %s: %s", userID, err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	return
}
