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
	// ✅ STEP 1: Load .env FIRST, before anything else
	loadEnv()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	loadPages()
	loadAdminPages()

	if err := connectDB(); err != nil {
		log.Fatal("database: ", err)
	}
	defer db.Close()

	mux := http.NewServeMux()
	
	// --- Static Files ---
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// --- Main Pages ---
	mux.HandleFunc("GET /{$}", HomeHandler)
	mux.Handle("GET /home", http.RedirectHandler("/", http.StatusMovedPermanently))
	mux.HandleFunc("GET /about", AboutHandler)
	mux.HandleFunc("GET /courses", CoursesHandler)
	mux.HandleFunc("GET /for-schools", SchoolHandler)
	
	// --- Apply / Register ---
	mux.HandleFunc("GET /apply", ApplyHandler)
	mux.HandleFunc("POST /apply", ApplySubmitHandler)
	mux.Handle("GET /register", http.RedirectHandler("/apply", http.StatusMovedPermanently))
	
	// --- Testimonials ---
	mux.HandleFunc("GET /testimonials", TestimonialHandler)
	mux.HandleFunc("POST /testimonials", SubmitTestimonialHandler)
	
	// --- Blog ---
	mux.HandleFunc("GET /blog", BlogHandler)
	mux.HandleFunc("GET /blog/{slug}", BlogPostHandler)
	
	// --- Contact ---
	mux.HandleFunc("GET /contact", ContactHandler)
	mux.HandleFunc("POST /contact", ContactSubmitHandler)
	
	// --- Subscribers ---
	mux.HandleFunc("POST /subscribe", SubscribeHandler)
	
	// --- ADMIN PORTAL ---
	mux.HandleFunc("GET /admin/login", AdminLoginHandler)
	mux.HandleFunc("POST /admin/login", AdminLoginSubmitHandler)
	mux.HandleFunc("GET /admin/logout", AdminLogoutHandler)
	
	mux.HandleFunc("GET /admin", AdminDashboardHandler)
	mux.HandleFunc("GET /admin/applications", AdminApplicationsHandler)
	mux.HandleFunc("GET /admin/messages", AdminMessagesHandler)
	mux.HandleFunc("GET /admin/subscribers", AdminSubscribersHandler)
	mux.HandleFunc("GET /admin/testimonials", AdminTestimonialsHandler)
	mux.HandleFunc("POST /admin/testimonials/toggle", AdminToggleTestimonialHandler)
	
	// --- Health Check ---
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