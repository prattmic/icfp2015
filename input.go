package main

import (
	"io"
	"encoding/json"
)

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
