package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/willvelida/go-rest-api/pkg/mocks"
)

func DeleteBook(w http.ResponseWriter, r *http.Request) {
	// Read the dynamic id parameter
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	// Iterate over all the mock books
	for index, book := range mocks.Books {
		if book.Id == id {
			// Delete and send responses when Ids match
			mocks.Books = append(mocks.Books[:index], mocks.Books[index+1:]...)
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			json.NewEncoder(w).Encode("Deleted")
			break
		}
	}
}
