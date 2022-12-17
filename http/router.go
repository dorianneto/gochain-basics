package http

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

func router() http.Handler {
	muxRouter := mux.NewRouter()

	muxRouter.HandleFunc("/", handleGetBlockchain).Methods("GET")
	muxRouter.HandleFunc("/", handleWriteBlock).Methods("POST")

	return muxRouter
}

func Run() error {
	router := router()

	port := os.Getenv("PORT")
	log.Println("Listening on ", port)

	server := &http.Server{
		Addr:           ":" + port,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := server.ListenAndServe(); err != nil {
		return err
	}

	return nil
}
