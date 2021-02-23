package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/3c1u/catechdojo/db"
)

// NOTE: 一般的にはユーザーとセッションを分けるが，一つのユーザーに対して一つのトークンしか結びつかないため
// ここでは考えないことにする．

type CharacterListResponse struct {
	Characters []UserCharacter `json:"characters"`
}

type UserCharacter struct {
	UserCharacterID string `json:"userCharacterID"`
	CharacterID     string `json:"characterID"`
	Name            string `json:"name"`
}

func HandleCharacterList(w http.ResponseWriter, r *http.Request) {
	var response CharacterListResponse

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

	c, err := db.EnumerateUserCharacters(user.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			map[string]string{
				"error":       "failed to enumerate a character",
				"description": err.Error(),
			},
		)
		return
	}

	userCharacters := []UserCharacter{}
	for i := 0; i < len(c); i++ {
		item := c[i]
		userCharacters = append(userCharacters, UserCharacter{
			UserCharacterID: fmt.Sprint(item.ID),
			CharacterID:     fmt.Sprint(item.CharacterID),
			Name:            item.Character.Name,
		})
	}

	response.Characters = userCharacters

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
