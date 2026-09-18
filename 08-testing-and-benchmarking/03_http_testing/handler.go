package httptesting

import (
	"encoding/json"
	"net/http"
)

type GreetResponse struct {
	Greeting string `json:"greeting"`
}

func GreetHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "Guest"
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(GreetResponse{
		Greeting: "Hello, " + name + "!",
	})
}
