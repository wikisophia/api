package arguments

import (
	"errors"
	"fmt"
)

// Argument is the core data type for the API.
type Argument struct {
	ID         int64    `json:"id"`
	Version    int      `json:"version"`
	Conclusion string   `json:"conclusion"`
	Premises   []string `json:"premises"`
}

// Validate returns nil if the argument is well-formed, or an error if not.
func (a *Argument) Validate() error {
	if a.Conclusion == "" {
		return errors.New("arguments must have a conclusion")
	}
	if len(a.Premises) < 2 {
		return errors.New("arguments must have at least 2 premises")
	}
	for i, premise := range a.Premises {
		if premise == "" {
			return fmt.Errorf("argument premise[%d] is empty, but must not be", i)
		}
	}

	return nil
}

// ByID can be used to sort slices of Arguments by ID.
type ByID []Argument

func (c ByID) Len() int {
	return len(c)
}
func (c ByID) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
func (c ByID) Less(i, j int) bool {
	return c[i].ID < c[j].ID
}
