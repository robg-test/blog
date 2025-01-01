package main

import (
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/a-h/templ"
	"github.com/gorilla/mux"
	"github.com/robgtest/blog/internal"
	"github.com/robgtest/blog/web/blogs"
	"github.com/robgtest/blog/web/blogs/stoicism"
	"github.com/robgtest/blog/web/pages"
)

func main() {
	loadableImages, err := getImagesList()
	if err != nil {
		panic(err)
	}
	log.SetFlags(log.LstdFlags)
	setup()

	router := mux.NewRouter()
	setupStaticHandlers(router, loadableImages)
	setupPageHandlers(router)
	setupBlogHandler(router)
	log.Println("Listening on :8080")
	http.ListenAndServe(":8080", internal.UserSessionManager.LoadAndSave(router))
}

func setupPageHandlers(router *mux.Router) {
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		indexPage := pages.IndexPage()
		templ.Handler(indexPage).ServeHTTP(w, r)
	})
}

func setupStaticHandlers(router *mux.Router, loadableImages []string) {
	router.HandleFunc("/styles.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/static/css/output.css")
	})
	router.HandleFunc("/images/{path}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		for key, value := range vars {
			if key == "path" {
				log.Printf("Loading image: %s", value)
				for _, image := range loadableImages {
					log.Printf("Lookup on Image: %s", image)
					if strings.Compare(value, image) == 0 {
						// process each image
						log.Println("Presenting Image")
						http.ServeFile(w, r, fmt.Sprintf("./web/static/images/%s", image))
					}
				}
			}
		}
	})
}

func getImagesList() ([]string, error) {
	directory := "./web/static/images"
	fileMap := []string{}

	err := filepath.WalkDir(directory, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			fileMap = append(fileMap, d.Name())
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return fileMap, nil
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

func setupBlogHandler(router *mux.Router) {
	router.HandleFunc("/blog/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		log.Println("Blog Requested")
		var blog templ.Component
		if vars["id"] == "1" {
			blog = blogs.BlogIntro()
		}
		if vars["id"] == "2" {
			blog = blogs.AWSServerlessBlog()
		}
		if vars["id"] == "3" {
			blog = stoicism.ControlAndChoice()
		}
		templ.Handler(blog).ServeHTTP(w, r)
	})
}
