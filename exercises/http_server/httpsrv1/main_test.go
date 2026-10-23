// httpsrv1
// Write an http.HandlerFunc and test it with httptest, no real network needed.

package main_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// greetHandler responds with "Hello, <name>!", where name comes from the
// "name" query parameter and defaults to "world".
func greetHandler(w http.ResponseWriter, r *http.Request) {
	// FIXME: read the name query param (default "world") and write the greeting.
	queryParams := r.URL.Query()
	qpNameValue := queryParams.Get("name")
	if qpNameValue == "" {
		qpNameValue = "world"
	}

	// Obtain the name from the `name` query parameter
	// Write response with that name and default
	payload := fmt.Sprintf("Hello, %s!", qpNameValue)
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(payload))
	if err != nil {
		panic(err)
	}
}

func TestGreetHandlerDefault(t *testing.T) {
	rec := httptest.NewRecorder()
	greetHandler(rec, httptest.NewRequest(http.MethodGet, "/greet", nil))
	if body, _ := io.ReadAll(rec.Result().Body); string(body) != "Hello, world!" {
		t.Errorf("want %q, got %q", "Hello, world!", string(body))
	}
}

func TestGreetHandlerWithName(t *testing.T) {
	rec := httptest.NewRecorder()
	greetHandler(rec, httptest.NewRequest(http.MethodGet, "/greet?name=Go", nil))
	if body, _ := io.ReadAll(rec.Result().Body); string(body) != "Hello, Go!" {
		t.Errorf("want %q, got %q", "Hello, Go!", string(body))
	}
}
