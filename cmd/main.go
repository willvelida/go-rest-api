package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/willvelida/go-rest-api/pkg/handlers"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/books", handlers.GetAllBooks).Methods(http.MethodGet)
	router.HandleFunc("/books", handlers.AddBook).Methods(http.MethodPost)
	router.HandleFunc("/books/{id}", handlers.GetBook).Methods(http.MethodGet)
	router.HandleFunc("/books/{id}", handlers.UpdateBook).Methods(http.MethodPut)
	router.HandleFunc("/books/{id}", handlers.DeleteBook).Methods(http.MethodDelete)
	router.HandleFunc("/health", handlers.HealthCheckHandler)

	log.Println("API is running")
	http.ListenAndServe(":8080", router)
}
