// httpsrv2
// Route requests with http.ServeMux using Go 1.22 method + path patterns.

package main_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// userHandler responds with "user <id>", taking id from the path —
// e.g. GET /users/42 -> "user 42".
func userHandler(w http.ResponseWriter, r *http.Request) {
	// FIXME: read the {id} path value and write it back.
	userId := r.PathValue("id")
	payload := fmt.Sprintf("user %s", userId)
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(payload))
	if err != nil {
		panic(err)
	}
}

// newMux wires the routes.
func newMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /users/{id}", userHandler)
	return mux
}

func TestUserRoute(t *testing.T) {
	rec := httptest.NewRecorder()
	newMux().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/users/42", nil))
	res := rec.Result()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("want 200, got %d", res.StatusCode)
	}
	if body, _ := io.ReadAll(res.Body); string(body) != "user 42" {
		t.Errorf("want %q, got %q", "user 42", string(body))
	}
}

func TestWrongMethodIsRejected(t *testing.T) {
	rec := httptest.NewRecorder()
	newMux().ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/users/42", nil))
	if got := rec.Result().StatusCode; got != http.StatusMethodNotAllowed {
		t.Errorf("want 405 for POST, got %d", got)
	}
}
