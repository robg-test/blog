package main

import (
	"encoding/xml"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/robgtest/blog/internal"
	"github.com/robgtest/blog/internal/models"
	"github.com/robgtest/blog/internal/static"
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

	router := chi.NewRouter()
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
		host := ":7000"
		log.Println("Starting development server on", host)
		err := http.ListenAndServe(host, internal.UserSessionManager.LoadAndSave(router))
		if err != nil {
			log.Fatalf("server failed: %s", err)
		}
	}
}

func setupPageHandlers(router *chi.Mux) {
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
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

	router.Put("/theme", func(w http.ResponseWriter, r *http.Request) {
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

func setupStaticHandlers(router *chi.Mux, loadableImages []string) {
	router.Get("/styles.css", serveCSS)
	router.Get("/prism.css", servePrismCSS)
	router.Get("/js/prism.js", servePrismJS)
	router.Get("/rss.xml", serveRSS)
	router.Get("/images/*", func(w http.ResponseWriter, r *http.Request) {
		serveImage(w, r, loadableImages)
	})
}

type rssItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

type rssChannel struct {
	Title       string    `xml:"title"`
	Link        string    `xml:"link"`
	Description string    `xml:"description"`
	Items       []rssItem `xml:"item"`
}

type rssFeed struct {
	XMLName xml.Name   `xml:"rss"`
	Version string     `xml:"version,attr"`
	Channel rssChannel `xml:"channel"`
}

func serveRSS(w http.ResponseWriter, r *http.Request) {
	posts := []models.BlogMeta{
		static.QuietSkillsData,
		static.GrugAutomationData,
		static.PerformanceWorkshop,
		static.IsCopilotADudData,
		static.AWSServerlessData,
		static.ToBeSteadyData,
		static.ControlAndChoiceData,
		static.IntroData,
	}

	items := make([]rssItem, len(posts))
	for i, p := range posts {
		items[i] = rssItem{
			Title:       p.Title,
			Link:        p.Url,
			Description: p.Description,
			PubDate:     p.Published.Format("Mon, 02 Jan 2006 00:00:00 +0000"),
		}
	}

	feed := rssFeed{
		Version: "2.0",
		Channel: rssChannel{
			Title:       "Bob Productions",
			Link:        "https://blog.bob-productions.dev/",
			Description: "Engineering, performance, and the occasional opinion.",
			Items:       items,
		},
	}

	w.Header().Set("Content-Type", "application/rss+xml; charset=utf-8")
	w.Write([]byte(xml.Header))
	xml.NewEncoder(w).Encode(feed)
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
	path := r.URL.Path[len("/images/"):]
	log.Println("Image path:", path)
	log.Println(loadableImages)
	if contains(loadableImages, "web/static/images/"+path) {
		log.Println("Serving image:", path)
		http.ServeFile(w, r, fmt.Sprintf("./web/static/images/%s", path))
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

func setupBlogHandler(router *chi.Mux) {
	router.Get("/blog/{id}", func(w http.ResponseWriter, r *http.Request) {
		theme := internal.GetMessage("theme", r)
		if theme == "" {
			theme = "retro"
		}

		log.Printf("Got theme: %v", theme)
		log.Println("Blog Requested")
		id := chi.URLParam(r, "id")

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
		case "quiet-skills":
			blog = blogs.QuietSkillsBlog(theme)
		}

		templ.Handler(blog).ServeHTTP(w, r)
	})
}
