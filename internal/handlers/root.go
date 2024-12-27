package handlers

import (
	"log"
	"net/http"

	"github.com/a-h/templ"
	"github.com/robgtest/blog/web/pages"
)

func LaunchIndexPage(r *http.Request, w http.ResponseWriter) {
	indexPage := pages.IndexPage()
	templ.Handler(indexPage).ServeHTTP(w, r)
	log.Println("Successfully served index Page")
}
