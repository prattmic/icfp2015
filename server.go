package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
)

// All active games are controlled by an AI.
var active map[string]AI = make(map[string]AI)

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

	// Ignore all but the first seeded game.
	g := GamesFromProblem(&problem)[0]
	active[k] = NewAI(g)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	log.Println("Game %s", g.B)

	if err := json.NewEncoder(w).Encode(g.B); err != nil {
		http.Error(w, fmt.Sprintf("Unable to encode JSON: %v", err), http.StatusInternalServerError)
		return
	}
}

func getFrames(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GET params:", r.URL.Query())

	token := r.URL.Query().Get("token")
	log.Printf("Token: %v", token)

	w.Header().Set("Server", "A Go Web Server")
	w.WriteHeader(200)
}

func runServer() {
	http.HandleFunc("/newgame", newGameHandler)
	http.HandleFunc("/getframes", getFrames)

	http.ListenAndServe(":8080", nil)
}
