package handler

import (
	"fmt"
	"net/http"
)

func HandleGachaDraw(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Hello, world!")
}
