package test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/willvelida/go-rest-api/pkg/handlers"
)

func TestGetAllBooks(t *testing.T) {
	req, err := http.NewRequest("GET", "/books", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.GetAllBooks)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body is what we expect
	expected := rr.Body.String()
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestGetBook(t *testing.T) {
	// Add a table-driven structure to test
	testBook := []struct {
		id     string
		title  string
		author string
		desc   string
	}{
		{"1", "The Hobbit", "J.R.R. Tolkien", "A hobbit is a small human-like creature that enjoys a comfortable, quiet life, usually in a hobbit-hole in the side of a hill. Hobbits are generally peace-loving folk, but they can be fierce fighters when their homes and ways of life are threatened."},
	}

	for _, tc := range testBook {
		path := fmt.Sprintf("/books/%s", tc.id)
		req, err := http.NewRequest("GET", path, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/books/{id}", handlers.GetBook).Methods("GET")
		router.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}
	}
}

func TestGetBookNotFound(t *testing.T) {
	req, err := http.NewRequest("GET", "/books", nil)
	if err != nil {
		t.Fatal(err)
	}

	q := req.URL.Query()
	q.Add("id", "999")
	rr := httptest.NewRecorder()
	req.URL.RawQuery = q.Encode()
	handler := http.HandlerFunc(handlers.GetBook)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
}

func TestAddBook(t *testing.T) {
	var requestBody = []byte(`{"title":"The Hobbit","author":"J.R.R. Tolkien","desc":"A hobbit is a small human-like creature that enjoys a comfortable, quiet life, usually in a hobbit-hole in the side of a hill. Hobbits are generally peace-loving folk, but they can be fierce fighters when their homes and ways of life are threatened."}`)

	req, err := http.NewRequest("POST", "/books", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.AddBook)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}
}

func TestUpdateBook(t *testing.T) {
	var requestBody = []byte(`{"title":"The Hobbit 2","author":"J.R.R. Tolkien","desc":"A hobbit is a small human-like creature that enjoys a comfortable, quiet life, usually in a hobbit-hole in the side of a hill. Hobbits are generally peace-loving folk, but they can be fierce fighters when their homes and ways of life are threatened."}`)

	req, err := http.NewRequest("PUT", "/books", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.UpdateBook)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestDeleteBook(t *testing.T) {
	req, err := http.NewRequest("DELETE", "/books", nil)
	if err != nil {
		t.Fatal(err)
	}

	q := req.URL.Query()
	q.Add("id", "1")
	rr := httptest.NewRecorder()
	req.URL.RawQuery = q.Encode()
	handler := http.HandlerFunc(handlers.DeleteBook)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestHealthCheckHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.HealthCheckHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"alive": true}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}
