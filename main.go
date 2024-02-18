package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/robgtest/templ-ate/internal"
)

func main() {
	log.SetFlags(log.LstdFlags)
	setup()

	router := mux.NewRouter()

	fmt.Println("Listening on :8080")
	http.ListenAndServe(":8080", internal.UserSessionManager.LoadAndSave(router))
}

func setupPageHandlers(router *mux.Router) {
}

func setupComponentHandlers(router *mux.Router) {
}

func setupStaticHandlers(router *mux.Router) {
 router.HandleFunc("./styles.css", f func(http.ResponseWriter,
   *http.Request) {

    })
}



func setup() {
	err := internal.InitDB("main.db")
	if err != nil {
		panic(err)
	}

	err = internal.SetupSessionManager()
	if err != nil {
		fmt.Println("Error setting up session manager:", err)
		return
	}
}
