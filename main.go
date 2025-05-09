package main

import (
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/a-h/templ"
	"github.com/gorilla/mux"
	"github.com/robgtest/blog/internal"
	"github.com/robgtest/blog/web/blogs"
	"github.com/robgtest/blog/web/blogs/stoicism"
	"github.com/robgtest/blog/web/pages"
)

const (
	defaultTheme   = "retro"
	synthwaveTheme = "synthwave"
)

func main() {
	loadableImages, err := getImagesList()
	if err != nil {
		log.Fatalf("Failed to get images list: %v", err)
	}
	log.SetFlags(log.LstdFlags)
	setup()

	router := mux.NewRouter()
	setupStaticHandlers(router, loadableImages)
	setupBlogHandler(router)
	setupPageHandlers(router)

	env := os.Getenv("ENV")

	if env == "production" {
		host := ":443"
		certPath := os.Getenv("CERT_PATH")
		keyPath := os.Getenv("KEY_PATH")

		log.Println("Starting secure server on", host)
		err := http.ListenAndServeTLS(host, certPath, keyPath, internal.UserSessionManager.LoadAndSave(router))
		if err != nil {
			log.Printf("secure server failed: %s", err)
		}
	} else {
		host := ":8080"
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
		if theme == "" {
			theme = defaultTheme
		}
		log.Printf("Got theme: %v", theme)
		log.Println("Blog Requested")
		indexPage := pages.IndexPage(theme)
		if indexPage == nil {
			indexPage = pages.IndexPage(defaultTheme)
		}
		templ.Handler(indexPage).ServeHTTP(w, r)
	})

	router.HandleFunc("/theme", func(w http.ResponseWriter, r *http.Request) {
		log.Println("change theme - old theme " + internal.GetMessage("theme", r))
		if internal.GetMessage("theme", r) != synthwaveTheme {
			internal.PutMessage("theme", synthwaveTheme, r)
			log.Println("change theme - new theme synthwave")
		} else {
			internal.PutMessage("theme", defaultTheme, r)
			log.Println("change theme - new theme retro")
		}
		w.WriteHeader(http.StatusOK)
	})
}

func setupStaticHandlers(router *mux.Router, loadableImages []string) {
	router.HandleFunc("/styles.css", serveCSS)
	router.HandleFunc("/prism.css", servePrismCSS)
	router.HandleFunc("/js/prism.js", servePrismJS)
	router.HandleFunc("/images/{path:.*}", func(w http.ResponseWriter, r *http.Request) {
		serveImage(w, r, loadableImages)
	})
}

func serveCSS(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./web/static/css/output.css")
}

func servePrismCSS(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./web/static/css/prism.css")
}

func servePrismJS(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./web/static/js/prism.js")
}

func serveImage(w http.ResponseWriter, r *http.Request, loadableImages []string) {
	vars := mux.Vars(r)
	log.Println(vars)
	log.Println(loadableImages)
	if contains(loadableImages, "web/static/images/"+vars["path"]) {
		log.Println(vars["path"])
		http.ServeFile(w, r, fmt.Sprintf("./web/static/images/%s", vars["path"]))
	}
}

func contains(images []string, path string) bool {
	for _, img := range images {
		if img == path {
			return true
		}
	}
	return false
}

func getImagesList() ([]string, error) {
	directory := "./web/static/images"
	fileMap := []string{}

	err := filepath.WalkDir(directory, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			fileMap = append(fileMap, path) // Use path to include subdirectory structure
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return fileMap, nil
}

func setup() {
	err := internal.InitDB(os.Getenv("TURSO_DATABASE"))
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
		if theme == "" {
			theme = "retro"
		}
		log.Printf("Got theme: %v", theme)
		log.Println("Blog Requested")
		id := vars["id"]

		var blog templ.Component
		switch id {
		case "intro":
			blog = blogs.BlogIntro(theme)
		case "serverless":
			blog = blogs.AWSServerlessBlog(theme)
		case "control-and-choice":
			blog = stoicism.ControlAndChoice(theme)
		case "to-be-steady":
			blog = stoicism.ToBeSteady(theme)
		case "ai-autocomplete":
			blog = blogs.IsCopilotADud(theme)
		case "perf-workshop":
			blog = blogs.PerformanceWorkshop(theme)
		}

		templ.Handler(blog).ServeHTTP(w, r)
	})
}
