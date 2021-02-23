package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/3c1u/catechdojo/db"
)

type GachaDrawRequest struct {
	Times int `json:"times"`
}

type GachaDrawResponse struct {
	Results []GachaResult `json:"results"`
}

type GachaResult struct {
	CharacterID string `json:"characterID"`
	Name        string `json:"Name"`
}

func HandleGachaDraw(w http.ResponseWriter, r *http.Request) {
	var request GachaDrawRequest
	var response GachaDrawResponse

	w.Header().Set("Content-Type", "application/json")

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(
			map[string]string{
				"error":       "failed to read the request",
				"description": err.Error(),
			},
		)
		return
	}

	if err = json.Unmarshal(reqBody, &request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(
			map[string]string{
				"error":       "failed to unmarshal request",
				"description": err.Error(),
			},
		)
		return
	}

	token := r.Header.Get("x-token")
	if token == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(
			map[string]string{
				"error":       "empty token",
				"description": "",
			},
		)
		return
	}

	user, err := db.UserGet(token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			map[string]string{
				"error":       "failed to get a user",
				"description": err.Error(),
			},
		)
		return
	}

	if user == nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(
			map[string]string{
				"error":       "failed to get a user",
				"description": "user not found",
			},
		)
		return
	}

	c, err := db.DrawGacha(user.ID, request.Times)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			map[string]string{
				"error":       "failed to draw gacha",
				"description": err.Error(),
			},
		)
		return
	}

	gachaResults := []GachaResult{}
	for i := 0; i < len(c); i++ {
		item := c[i]
		gachaResults = append(gachaResults, GachaResult{
			Name:        item.Character.Name,
			CharacterID: fmt.Sprint(item.CharacterID),
		})
	}

	response.Results = gachaResults

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
