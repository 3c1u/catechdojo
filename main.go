package main

import (
	"log"
	"net/http"
	"time"

	"github.com/3c1u/catechdojo/db"
	"github.com/3c1u/catechdojo/handler"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	log.Println("Connecting to database...")
	db.Init()

	log.Println("Launching server...")
	router := mux.NewRouter()
	router.HandleFunc("/", handler.HandleRoot)

	// user
	router.HandleFunc("/user/create", handler.HandleUserCreate).Methods("POST")
	router.HandleFunc("/user/get", handler.HandleUserGet).Methods("GET")
	router.HandleFunc("/user/update", handler.HandleUserUpdate).Methods("PUT")

	// gacha
	router.HandleFunc("/gacha/draw", handler.HandleGachaDraw).Methods("POST")

	// character
	router.HandleFunc("/character/list", handler.HandleCharacterList).Methods("GET")

	// NOTE: Swagger Editorでテストするために，CORS許可している．
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodOptions,
			http.MethodHead,
		},
		AllowedHeaders: []string{
			"*",
		},
		AllowCredentials: true,
	})
	handler := c.Handler(router)

	// TODO: 環境変数から読み取る
	addr := "0.0.0.0:8080"

	server := &http.Server{
		Handler:      handler,
		Addr:         addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("Listening to:", addr)
	log.Fatal(server.ListenAndServe())
}
