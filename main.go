package main

import (
	"log"
	"net/http"

	"github.com/a-h/templ"
	"github.com/gorilla/mux"
	"github.com/robgtest/templ-ate/internal"
	"github.com/robgtest/templ-ate/web/pages"
)

func main() {
	log.SetFlags(log.LstdFlags)
	setup()

	router := mux.NewRouter()
	setupStaticHandlers(router)
	setupPageHandlers(router)
	log.Println("Listening on :8080")
	http.ListenAndServe(":8080", internal.UserSessionManager.LoadAndSave(router))
}

func setupPageHandlers(router *mux.Router) {
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		indexPage := pages.IndexPage()
		templ.Handler(indexPage).ServeHTTP(w, r)
	})
}

func setupComponentHandlers(router *mux.Router) {
}

func setupStaticHandlers(router *mux.Router) {
	router.HandleFunc("/styles.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/static/css/output.css")
	})
}

func setup() {
	err := internal.InitDB("main.db")
	if err != nil {
		panic(err)
	}

	err = internal.SetupSessionManager()
	if err != nil {
		log.Println("Error setting up session manager:", err)
		return
	}
}
