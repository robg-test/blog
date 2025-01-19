package main

import (
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
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
	setupBlogHandler(router)
	setupPageHandlers(router)

	env := os.Getenv("ENV")
	host := ":8080"

	if env == "production" {
		host = ":443"
		certPath := os.Getenv("CERT_PATH")
		keyPath := os.Getenv("KEY_PATH")

		log.Println("Starting secure server on", host)
		err := http.ListenAndServeTLS(host, certPath, keyPath, internal.UserSessionManager.LoadAndSave(router))
		if err != nil {
			log.Printf("secure server failed: %s", err)
		}
	} else {
		log.Println("Starting development server on", host)
		err := http.ListenAndServe(host, internal.UserSessionManager.LoadAndSave(router))
		if err != nil {
			log.Fatalf("server failed: %s", err)
		}
	}
}

func setupPageHandlers(router *mux.Router) {
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		theme := internal.GetMessage("theme", r)
		log.Printf("Got theme: %v", theme)
		log.Println("Blog Requested")
		indexPage := pages.IndexPage(theme)
		templ.Handler(indexPage).ServeHTTP(w, r)
	})

	router.HandleFunc("/theme", func(w http.ResponseWriter, r *http.Request) {
		log.Println("change theme - old theme " + internal.GetMessage("theme", r))
		if internal.GetMessage("theme", r) != "synthwave" {
			internal.PutMessage("theme", "synthwave", r)
			log.Println("change theme - new theme synthwave")
		} else {
			internal.PutMessage("theme", "retro", r)
			log.Println("change theme - new theme retro")
		}
		w.WriteHeader(http.StatusOK)
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
		theme := internal.GetMessage("theme", r)
		log.Printf("Got theme: %v", theme)
		log.Println("Blog Requested")
		var blog templ.Component
		switch vars["id"] {
		case "1":
			blog = blogs.BlogIntro(theme)
		case "2":
			blog = blogs.AWSServerlessBlog(theme)
		case "3":
			blog = stoicism.ControlAndChoice(theme)
		case "4":
			blog = stoicism.ToBeSteady(theme)
		case "5":
			blog = blogs.IsCopilotADud(theme)

		}
		templ.Handler(blog).ServeHTTP(w, r)
	})
}
