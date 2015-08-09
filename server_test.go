package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestNewGame(t *testing.T) {
	record := httptest.NewRecorder()

	d, err := os.Open("qualifiers/problem_0.json")
	if err != nil {
		t.Fatalf("ioutil.ReadFile err: got %v want nil", err)
	}

	req := &http.Request{
		Method: "POST",
		Body:   d,
	}
	newGameHandler(record, req)

	if record.Code != 201 {
		t.Errorf("record.Code got %d want 201", record.Code)
	}

	var resp newGameResponse
	if err := json.NewDecoder(record.Body).Decode(&resp); err != nil {
		t.Errorf("Decode(%s) err: got %v want nil", record.Body, err)
	}
	// if resp.Token == "" {
	// 	t.Errorf(`Token: got "", want something`)
	// }
}
