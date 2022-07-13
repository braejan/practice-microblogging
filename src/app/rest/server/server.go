package server

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Binder func(server Server, router *mux.Router)
type Config struct {
	Port string
}

type Server interface {
	Config() *Config
}

func NewServer(ctx context.Context, configuration *Config) (broker *Broker, err error) {
	broker = &Broker{
		config: configuration,
		router: mux.NewRouter(),
	}
	return
}

type Broker struct {
	config *Config
	router *mux.Router
}

func (b *Broker) Config() (configuration *Config) {
	return b.config
}

func (broker *Broker) Start(binder Binder) {
	broker.router = mux.NewRouter()
	binder(broker, broker.router)
	log.Printf("Starting server on port %s", broker.config.Port)
	err := http.ListenAndServe(fmt.Sprintf(":%s", broker.config.Port), broker.router)
	if err != nil {
		log.Fatalf("error starting server: %v", err)
	}
}
