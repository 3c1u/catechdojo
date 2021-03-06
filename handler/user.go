package handler

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/3c1u/catechdojo/db"
)

type UserCreateRequest struct {
	Name string `json:"name"`
}

type UserUpdateRequest struct {
	Name string `json:"name"`
}

type UserCreateResponse struct {
	Token string `json:"token"`
}

type UserGetResponse struct {
	Name string `json:"name"`
}

func HandleUserCreate(w http.ResponseWriter, r *http.Request) {
	var request UserCreateRequest
	var response UserCreateResponse

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

	user, err := db.UserCreate(request.Name)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			map[string]string{
				"error":       "failed to create a user",
				"description": err.Error(),
			},
		)
		return
	}

	response.Token = user.Token

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func HandleUserGet(w http.ResponseWriter, r *http.Request) {
	var response UserGetResponse

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
		log.Println(err)
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
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(
			map[string]string{
				"error":       "failed to get a user",
				"description": "user not found",
			},
		)
		return
	}

	response.Name = user.Name

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func HandleUserUpdate(w http.ResponseWriter, r *http.Request) {
	var request UserUpdateRequest

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
		log.Println(err)
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
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(
			map[string]string{
				"error":       "failed to get a user",
				"description": "user not found",
			},
		)
		return
	}

	user.Name = request.Name
	err = db.UserUpdate(user)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			map[string]string{
				"error":       "failed to update a user profile",
				"description": err.Error(),
			},
		)
		return
	}

	w.WriteHeader(http.StatusOK)
}
