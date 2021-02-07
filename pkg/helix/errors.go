package helix

import "fmt"

type helixError struct {
	Type    string `json:"error"`
	Status  int
	Message string
}

func (h *helixError) Error() string {
	return fmt.Sprintf("%v %s: %s", h.Status, h.Type, h.Message)
}
