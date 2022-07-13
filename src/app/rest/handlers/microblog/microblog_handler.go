package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/braejan/practice-microblogging/src/app/rest/server"
	commonUsecases "github.com/braejan/practice-microblogging/src/domain/common/usecases"
	microblogUsecases "github.com/braejan/practice-microblogging/src/domain/microblog/usecases"
	userUsecases "github.com/braejan/practice-microblogging/src/domain/user/usecases"
	"github.com/gorilla/mux"
)

type allUsecases struct {
	userUC      *userUsecases.UserUsecases
	commonUC    *commonUsecases.CommonUsecases
	microblogUC *microblogUsecases.MicroblogUsecases
}

func newUsecases() (uc *allUsecases, err error) {
	uc = &allUsecases{}
	uc.userUC, err = userUsecases.NewUserUsecases()
	if err != nil {
		log.Printf("error creating userUsecases: %s", err.Error())
		return nil, err
	}
	uc.commonUC = commonUsecases.NewCommonUsecases()
	uc.microblogUC, err = microblogUsecases.NewMicroblogUsecases(uc.userUC, uc.commonUC)
	if err != nil {
		log.Printf("error creating userUsecases: %s", err.Error())
		return nil, err
	}
	return
}

func GetAllPostHandler(server server.Server) (handlerFunc http.HandlerFunc) {
	handlerFunc = func(w http.ResponseWriter, r *http.Request) {
		requestUserID := r.Header.Get("user")
		if requestUserID == "" {
			w.WriteHeader(http.StatusBadRequest)
			log.Printf("error getting userID from header")
			return
		}
		uc, err := newUsecases()
		if err != nil {
			log.Printf("error init uc: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		posts, err := uc.microblogUC.GetAllPosts()
		if err != nil {
			log.Printf("error getting posts: %s", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if len(posts) == 0 {
			log.Println("posts not found")
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(posts)
		if err != nil {
			log.Printf("error writing posts: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	return
}

type newPost struct {
	Text string `json:"text"`
}

func CreatePostHandler(server server.Server) (handlerFunc http.HandlerFunc) {
	handlerFunc = func(w http.ResponseWriter, r *http.Request) {
		requestUserID := r.Header.Get("user")
		if requestUserID == "" {
			w.WriteHeader(http.StatusBadRequest)
			log.Printf("error getting userID from header")
			return
		}
		uc, err := newUsecases()
		if err != nil {
			log.Printf("error init uc: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		var requestText newPost
		err = json.NewDecoder(r.Body).Decode(&requestText)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Printf("error getting text from body")
			return
		}
		err = uc.microblogUC.CreatePost(requestUserID, requestText.Text)
		if err != nil {
			log.Printf("error creating post: %s", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
	}
	return
}

func GetPostByIDHandler(server server.Server) (handlerFunc http.HandlerFunc) {
	handlerFunc = func(w http.ResponseWriter, r *http.Request) {
		postID := mux.Vars(r)["postID"]
		if postID == "" {
			w.WriteHeader(http.StatusBadRequest)
			log.Printf("error getting postID from path")
			return
		}
		requestUserID := r.Header.Get("user")
		if requestUserID == "" {
			w.WriteHeader(http.StatusBadRequest)
			log.Printf("error getting userID from header")
			return
		}
		uc, err := newUsecases()
		if err != nil {
			log.Printf("error init uc: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		post, err := uc.microblogUC.GetPostByID(postID, requestUserID, false)
		if err != nil {
			log.Printf("error getting posts: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if post == nil {
			log.Println("posts not found")
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(post)
		if err != nil {
			log.Printf("error writing posts: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	return
}

func GetAllPostByUserIDHandler(server server.Server) (handlerFunc http.HandlerFunc) {
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
		uc, err := newUsecases()
		if err != nil {
			log.Printf("error init uc: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		post, err := uc.microblogUC.GetAllPostsByUserID(userID)
		if err != nil {
			log.Printf("error getting posts: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if post == nil {
			log.Println("posts not found")
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(post)
		if err != nil {
			log.Printf("error writing posts: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	return
}

func LikePostHandler(server server.Server) (handlerFunc http.HandlerFunc) {
	handlerFunc = func(w http.ResponseWriter, r *http.Request) {
		requestUserID := r.Header.Get("user")
		if requestUserID == "" {
			w.WriteHeader(http.StatusBadRequest)
			log.Printf("error getting userID from header")
			return
		}
		postID := mux.Vars(r)["postID"]
		if postID == "" {
			w.WriteHeader(http.StatusBadRequest)
			log.Printf("error getting postID from path")
			return
		}
		uc, err := newUsecases()
		if err != nil {
			log.Printf("error init uc: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		err = uc.microblogUC.LikePost(postID, requestUserID)
		if err != nil {
			log.Printf("error like post: %s", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
	}
	return
}

func DislikePostHandler(server server.Server) (handlerFunc http.HandlerFunc) {
	handlerFunc = func(w http.ResponseWriter, r *http.Request) {
		requestUserID := r.Header.Get("user")
		if requestUserID == "" {
			w.WriteHeader(http.StatusBadRequest)
			log.Printf("error getting userID from header")
			return
		}
		postID := mux.Vars(r)["postID"]
		if postID == "" {
			w.WriteHeader(http.StatusBadRequest)
			log.Printf("error getting postID from path")
			return
		}
		uc, err := newUsecases()
		if err != nil {
			log.Printf("error init uc: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		err = uc.microblogUC.DislikePost(postID, requestUserID)
		if err != nil {
			log.Printf("error like post: %s", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
	}
	return
}
