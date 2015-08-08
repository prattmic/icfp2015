package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
)

// Game is a fake game type, until mgyenik gives us a real one :)
type Game struct {
	problem InputProblem
}

func NewGame(p InputProblem) Game {
	return Game{
		problem: p,
	}
}

// All active games.
var games map[string]Game = make(map[string]Game)

type newGameResponse struct {
	// Token is used to reference your game in later requests.
	Token string
}

// POST a JSON InputProblem, receive a newGameResponse with the token to send
// to other methods.
func newGameHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	d := json.NewDecoder(r.Body)

	var problem InputProblem
	if err := d.Decode(&problem); err != nil {
		http.Error(w, fmt.Sprintf("Unable to parse JSON: %v", err), http.StatusBadRequest)
		return
	}

	// Generate a new key. Guaranteed random by fair dice roll!
	// TODO(prattmic): collision detection
	k := strconv.Itoa(rand.Int())
	games[k] = NewGame(problem)
	log.Printf("New game %s: %+v", k, games[k])

	resp := newGameResponse{Token: k}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, fmt.Sprintf("Unable to encode JSON: %v", err), http.StatusInternalServerError)
		return
	}
}

func runServer() {
	http.HandleFunc("/newgame", newGameHandler)

	http.ListenAndServe(":8080", nil)
}