package main

import (
	"encoding/gob"
	"github.com/marpio/webapp2/server/env"
	"github.com/marpio/webapp2/server/handlers"
	"github.com/marpio/webapp2/server/middleware"
	"github.com/marpio/webapp2/server/models"
	"github.com/marpio/webapp2/server/viewmodels"
	"github.com/stretchr/graceful"
	"net/http"
	"time"
)

func NewServer(middle http.Handler) *graceful.Server {
	srv := &graceful.Server{
		Timeout: 1 * time.Second,
		Server: &http.Server{
			Addr:    ":8000",
			Handler: middle,
		},
	}
	return srv
}

func main() {
	gob.Register(&viewmodels.OrderDto{})
	gob.Register(&models.UserRow{})

	environment := env.NewEnv()
	middle := middleware.NewDefaultMiddlewareChain(environment)
	router := handlers.NewConfiguredRouter(environment)
	h := middle.Then(router)
	srv := NewServer(h)
	srv.ListenAndServe()
}
