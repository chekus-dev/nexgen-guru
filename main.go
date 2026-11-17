package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func loadEnv() {
	b, err := os.ReadFile(".env")
	if err != nil {
		return // no .env file is fine (e.g. on Render)
	}
	for _, line := range strings.Split(string(b), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		k, v, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		v = strings.Trim(strings.TrimSpace(v), `"'`)
		if os.Getenv(strings.TrimSpace(k)) == "" {
			os.Setenv(strings.TrimSpace(k), v)
		}
	}
}
func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	loadPages()
	if err := connectDB(); err != nil {
		log.Fatal("database: ", err)
	}
	defer db.Close()

	mux := http.NewServeMux()
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	mux.HandleFunc("GET /{$}", HomeHandler)
	mux.Handle("GET /home", http.RedirectHandler("/", http.StatusMovedPermanently))
	mux.HandleFunc("GET /about", AboutHandler)
	mux.HandleFunc("GET /courses", CoursesHandler)
	mux.HandleFunc("GET /for-schools", SchoolHandler)
	mux.HandleFunc("GET /apply", ApplyHandler)
	mux.HandleFunc("POST /apply", ApplySubmitHandler)
	mux.Handle("GET /register", http.RedirectHandler("/apply", http.StatusMovedPermanently))
	mux.HandleFunc("GET /testimonials", TestimonialHandler)
	mux.HandleFunc("GET /blog", BlogHandler)
	mux.HandleFunc("GET /blog/{slug}", BlogPostHandler)
	mux.HandleFunc("GET /contact", ContactHandler)
	mux.HandleFunc("POST /contact", ContactSubmitHandler)
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("ok")) })

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	log.Printf("Starting the server on localhost:%s\n", port)
	log.Fatal(srv.ListenAndServe())
}
