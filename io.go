package main

import (
	"encoding/json"
	"io"
)

type InputProblem struct {
	Id           int
	Units        []Unit
	Width        int
	Height       int
	Filled       []Cell
	SourceLength int
	SourceSeeds  []uint64
}

type OutputEntry struct {
	ProblemId int    `json:"problemId"`
	Seed      uint64 `json:"seed"`
	Tag       string `json:"tag"`
	Solution  string `json:"solution"`
}

// This takes an io.Reader, and tries to unmarshal a JSON formatted
// InputProblem from it, returning an InputProblem to you.
func ParseInputProblem(r io.Reader) (*InputProblem, error) {
	var problem InputProblem

	d := json.NewDecoder(r)
	if err := d.Decode(&problem); err != nil {
		return nil, err
	}

	return &problem, nil
}
