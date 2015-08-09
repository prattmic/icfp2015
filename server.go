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
	a := NewAI(g)
	active[k] = a

	frames := []Board{}

	// All frames match the last currently, because I don't
	// understand how go handles references. Q.Q
	i := 1
	for {
		frame := *a.Game().B

		// Make this copy the object to save state for later.
		frames = append(frames, frame)
		done, err := a.Next()
		if done {
			log.Println("Game done!")
			break
		} else if err != nil {
			log.Printf("a.Next error: %v", err)
			break
		}
		i++
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(frames); err != nil {
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
