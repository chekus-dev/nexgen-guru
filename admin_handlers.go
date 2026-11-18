package main

import (
	"bytes" // <-- ADDED
	"crypto/rand"
	"encoding/base64"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

var adminSessions = map[string]time.Time{}

func generateToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

func isAuthenticated(r *http.Request) bool {
	cookie, err := r.Cookie("admin_session")
	if err != nil {
		return false
	}
	expiry, ok := adminSessions[cookie.Value]
	if !ok || time.Now().After(expiry) {
		delete(adminSessions, cookie.Value)
		return false
	}
	return true
}

func requireAuth(w http.ResponseWriter, r *http.Request) bool {
	if !isAuthenticated(r) {
		http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
		return false
	}
	return true
}

var adminPages = map[string]*template.Template{}

func loadAdminPages() {
	for _, n := range []string{"login", "dashboard", "applications", "messages", "subscribers", "testimonials"} {
		adminPages[n] = template.Must(template.ParseFiles("templates/admin_layout.html", "templates/admin_"+n+".html"))
	}
}

func renderAdmin(w http.ResponseWriter, name string, status int, d map[string]interface{}) {
	t, ok := adminPages[name]
	if !ok {
		http.Error(w, "page not found", http.StatusInternalServerError)
		return
	}
	if _, ok := d["LoggedIn"]; !ok {
		d["LoggedIn"] = false
	}
	
	// FIX: Declare the buffer before using it
	var buf bytes.Buffer
	if err := t.ExecuteTemplate(&buf, "admin_layout", d); err != nil {
		log.Println("admin render error:", err)
		http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	buf.WriteTo(w)
}

// ---------- Handlers ----------

func AdminLoginHandler(w http.ResponseWriter, r *http.Request) {
	if isAuthenticated(r) {
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}
	renderAdmin(w, "login", http.StatusOK, map[string]interface{}{"LoggedIn": false})
}

func AdminLoginSubmitHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	password := r.PostFormValue("password")
	expected := os.Getenv("ADMIN_PASSWORD")
	if expected == "" {
		expected = "nexgen2025"
	}
	
	if password != expected {
		renderAdmin(w, "login", http.StatusUnauthorized, map[string]interface{}{
			"LoggedIn": false,
			"Error":    "Invalid password. Please try again.",
		})
		return
	}

	token := generateToken()
	adminSessions[token] = time.Now().Add(24 * time.Hour)
	
	http.SetCookie(w, &http.Cookie{
		Name:     "admin_session",
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func AdminLogoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("admin_session")
	if err == nil {
		delete(adminSessions, cookie.Value)
	}
	http.SetCookie(w, &http.Cookie{
		Name:   "admin_session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
	http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
}

func AdminDashboardHandler(w http.ResponseWriter, r *http.Request) {
	if !requireAuth(w, r) { return }
	
	var appCount, msgCount, subCount, testCount int
	db.QueryRow("SELECT COUNT(*) FROM applications").Scan(&appCount)
	db.QueryRow("SELECT COUNT(*) FROM messages").Scan(&msgCount)
	db.QueryRow("SELECT COUNT(*) FROM subscribers").Scan(&subCount)
	db.QueryRow("SELECT COUNT(*) FROM testimonials").Scan(&testCount)
	
	renderAdmin(w, "dashboard", http.StatusOK, map[string]interface{}{
		"LoggedIn":  true,
		"AppCount":  appCount,
		"MsgCount":  msgCount,
		"SubCount":  subCount,
		"TestCount": testCount,
	})
}

func AdminApplicationsHandler(w http.ResponseWriter, r *http.Request) {
	if !requireAuth(w, r) { return }
	
	rows, err := db.Query(`
		SELECT id, 
		       COALESCE(name, '') as name, 
		       COALESCE(email, '') as email, 
		       COALESCE(phone, '') as phone, 
		       COALESCE(age, 0) as age, 
		       COALESCE(course, '') as course, 
		       COALESCE(location, '') as location, 
		       created_at 
		FROM applications 
		ORDER BY created_at DESC`)
	if err != nil {
		log.Println("AdminApplicationsHandler SQL Error:", err)
		renderAdmin(w, "applications", http.StatusOK, map[string]interface{}{
			"LoggedIn":     true,
			"Applications": []struct{}{},
			"DBError":      "Could not load applications: " + err.Error(),
		})
		return
	}
	defer rows.Close()
	
	type AppRow struct {
		ID        int
		Name      string
		Email     string
		Phone     string
		Age       int
		Course    string
		Location  string
		CreatedAt time.Time
	}
	var apps []AppRow
	for rows.Next() {
		var a AppRow
		if err := rows.Scan(&a.ID, &a.Name, &a.Email, &a.Phone, &a.Age, &a.Course, &a.Location, &a.CreatedAt); err != nil {
			log.Println("AdminApplicationsHandler Scan Error:", err)
			continue
		}
		apps = append(apps, a)
	}
	
	renderAdmin(w, "applications", http.StatusOK, map[string]interface{}{
		"LoggedIn":     true,
		"Applications": apps,
	})
}

func AdminMessagesHandler(w http.ResponseWriter, r *http.Request) {
	if !requireAuth(w, r) { return }
	
	rows, err := db.Query(`
		SELECT id, 
		       COALESCE(name, '') as name, 
		       COALESCE(email, '') as email, 
		       COALESCE(phone, '') as phone, 
		       COALESCE(interest, '') as interest, 
		       COALESCE(message, '') as message, 
		       created_at 
		FROM messages 
		ORDER BY created_at DESC`)
	if err != nil {
		log.Println("AdminMessagesHandler SQL Error:", err)
		renderAdmin(w, "messages", http.StatusOK, map[string]interface{}{
			"LoggedIn": true,
			"Messages": []struct{}{},
			"DBError":  "Could not load messages: " + err.Error(),
		})
		return
	}
	defer rows.Close()
	
	type MsgRow struct {
		ID        int
		Name      string
		Email     string
		Phone     string
		Interest  string
		Message   string
		CreatedAt time.Time
	}
	var msgs []MsgRow
	for rows.Next() {
		var m MsgRow
		if err := rows.Scan(&m.ID, &m.Name, &m.Email, &m.Phone, &m.Interest, &m.Message, &m.CreatedAt); err != nil {
			log.Println("AdminMessagesHandler Scan Error:", err)
			continue
		}
		msgs = append(msgs, m)
	}
	
	renderAdmin(w, "messages", http.StatusOK, map[string]interface{}{
		"LoggedIn": true,
		"Messages": msgs,
	})
}

func AdminSubscribersHandler(w http.ResponseWriter, r *http.Request) {
	if !requireAuth(w, r) { return }
	
	rows, err := db.Query(`
		SELECT id, 
		       COALESCE(email, '') as email, 
		       created_at 
		FROM subscribers 
		ORDER BY created_at DESC`)
	if err != nil {
		log.Println("AdminSubscribersHandler SQL Error:", err)
		renderAdmin(w, "subscribers", http.StatusOK, map[string]interface{}{
			"LoggedIn":    true,
			"Subscribers": []struct{}{},
			"DBError":     "Could not load subscribers: " + err.Error(),
		})
		return
	}
	defer rows.Close()
	
	type SubRow struct {
		ID        int
		Email     string
		CreatedAt time.Time
	}
	var subs []SubRow
	for rows.Next() {
		var s SubRow
		if err := rows.Scan(&s.ID, &s.Email, &s.CreatedAt); err != nil {
			log.Println("AdminSubscribersHandler Scan Error:", err)
			continue
		}
		subs = append(subs, s)
	}
	
	renderAdmin(w, "subscribers", http.StatusOK, map[string]interface{}{
		"LoggedIn":    true,
		"Subscribers": subs,
	})
}

func AdminTestimonialsHandler(w http.ResponseWriter, r *http.Request) {
	if !requireAuth(w, r) { return }
	
	rows, err := db.Query(`
		SELECT id, 
		       COALESCE(name, '') as name, 
		       COALESCE(email, '') as email, 
		       COALESCE(course, '') as course, 
		       COALESCE(quote, '') as quote, 
		       COALESCE(rating, 0) as rating, 
		       COALESCE(is_public, true) as is_public, 
		       created_at 
		FROM testimonials 
		ORDER BY created_at DESC`)
	if err != nil {
		log.Println("AdminTestimonialsHandler SQL Error:", err)
		renderAdmin(w, "testimonials", http.StatusOK, map[string]interface{}{
			"LoggedIn":     true,
			"Testimonials": []struct{}{},
			"DBError":      "Could not load testimonials: " + err.Error(),
		})
		return
	}
	defer rows.Close()
	
	type TestRow struct {
		ID        int
		Name      string
		Email     string
		Course    string
		Quote     string
		Rating    int
		IsPublic  bool
		CreatedAt time.Time
	}
	var tests []TestRow
	for rows.Next() {
		var t TestRow
		if err := rows.Scan(&t.ID, &t.Name, &t.Email, &t.Course, &t.Quote, &t.Rating, &t.IsPublic, &t.CreatedAt); err != nil {
			log.Println("AdminTestimonialsHandler Scan Error:", err)
			continue
		}
		tests = append(tests, t)
	}
	
	renderAdmin(w, "testimonials", http.StatusOK, map[string]interface{}{
		"LoggedIn":     true,
		"Testimonials": tests,
	})
}

func AdminToggleTestimonialHandler(w http.ResponseWriter, r *http.Request) {
	if !requireAuth(w, r) { return }
	
	r.ParseForm()
	id := r.PostFormValue("id")
	_, err := db.Exec(`UPDATE testimonials SET is_public = NOT is_public WHERE id = $1`, id)
	if err != nil {
		log.Println("toggle testimonial error:", err)
		http.Error(w, "Could not update testimonial: "+err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/admin/testimonials", http.StatusSeeOther)
}