package handler

import (
	"fmt"
	"net/http"
)

func HandleCharacterList(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Hello, world!")
}
