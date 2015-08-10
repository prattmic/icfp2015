package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Frame struct {
	BoardDelta []BoardCell
	Unit       *Unit
	Score      float64
	AI         string
}

type GameSolveResponse struct {
	Frames []Frame
	Board  *Board
}

type ReceivedProblem struct {
	Problem InputProblem
	AI      string
}

func getFrameDeltas(prev *Board, curr *Board) []BoardCell {
	deltas := []BoardCell{}

	if prev == nil {
		log.Println("prev was nil")
		return deltas
	}

	prev = prev.Fork()
	curr = curr.Fork()

	for y := range prev.Cells {
		for x := range prev.Cells[y] {
			if prev.Cells[x][y].Filled != curr.Cells[x][y].Filled {
				log.Printf("delta at %s, %s", x, y)
				deltas = append(deltas, curr.Cells[x][y])
			}
		}
	}

	return deltas
}

// POST a JSON InputProblem, receive a newGameResponse with the token to send
// to other methods.
func newGameHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("New request received!")
	if r.Method != "POST" {
		log.Printf("Not a post request")
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	d := json.NewDecoder(r.Body)

	var problem ReceivedProblem
	if err := d.Decode(&problem); err != nil {
		log.Printf("Error decoding! %v", err)
		http.Error(w, fmt.Sprintf("Unable to parse JSON: %v", err), http.StatusBadRequest)
		return
	}

	// Ignore all but the first seeded game.
	g := GamesFromProblem(&problem.Problem)[0]
	a := NewAI(g, problem.AI)

	response := GameSolveResponse{
		Board:  g.B.Fork(),
		Frames: []Frame{},
	}

	var prevBoard *Board

	i := 1
	for {
		game := a.Game()

		done, err := a.Next()

		deltas := getFrameDeltas(prevBoard, game.B)

		frame := Frame{
			BoardDelta: deltas,
			Unit:       game.currUnit.DeepCopy(),
			Score:      game.Score(),
			AI:         aiFlag,
		}

		prevBoard = game.B.Fork()

		response.Frames = append(response.Frames, frame)

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

	if err := json.NewEncoder(w).Encode(response); err != nil {
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
