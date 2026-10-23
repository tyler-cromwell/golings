// httpsrv3
// Decode a JSON request body and encode a JSON response with encoding/json.

package main_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type echoReq struct {
	Name string `json:"name"`
}

type echoResp struct {
	Greeting string `json:"greeting"`
}

// echoHandler decodes {"name": "..."} and responds with
// {"greeting": "Hello, <name>!"} and Content-Type application/json.
func echoHandler(w http.ResponseWriter, r *http.Request) {
	// FIXME: decode the JSON body, set the JSON content-type, then encode the reply.
	var payload echoReq
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte("Failed to read & decode request body"))
		if err != nil {
			panic(err)
		}
	}

	responseBody := echoResp{Greeting: fmt.Sprintf("Hello, %s!", payload.Name)}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(responseBody)
	if err != nil {
		panic(err)
	}
}

func TestEchoHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/echo", strings.NewReader(`{"name":"Go"}`))
	rec := httptest.NewRecorder()
	echoHandler(rec, req)

	res := rec.Result()
	if ct := res.Header.Get("Content-Type"); ct != "application/json" {
		t.Errorf("want Content-Type application/json, got %q", ct)
	}
	var got echoResp
	if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if got.Greeting != "Hello, Go!" {
		t.Errorf("want greeting %q, got %q", "Hello, Go!", got.Greeting)
	}
}
