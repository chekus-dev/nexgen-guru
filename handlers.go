package main

import (
	"bytes" // <--- ADDED THIS
	"html/template"
	"log"
	"net/http"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"
)

// ---------- data ----------

type Application struct {
	Name     string    `json:"name"`
	Phone    string    `json:"phone"`
	Age      int       `json:"age"`
	Course   string    `json:"course"`
	Location string    `json:"location"`
	Received time.Time `json:"received"`
}

type Post struct {
	Title, Slug, Body string
	Date              time.Time
}

var posts = []Post{} // add real posts here

type pageData struct {
	Title   string
	Site    SiteInfo
	Courses []Course
	Posts   []Post
	Post    *Post
	Form    map[string]string
	Errors  map[string]string
	Sent    bool
}

// ---------- templates ----------

var pages = map[string]*template.Template{}

// Parse once at startup, not on every request.
func loadPages() {
	for _, n := range []string{"home", "about", "courses", "schools", "apply", "testimonials", "blog", "post", "contact"} {
		pages[n] = template.Must(template.ParseFiles("templates/layout.html", "templates/"+n+".html"))
	}
}

// UPDATED render function to catch template errors
func render(w http.ResponseWriter, name string, status int, d pageData) {
	t, ok := pages[name]
	if !ok {
		http.Error(w, "page not found", http.StatusInternalServerError)
		return
	}
	d.Site = site
	d.Courses = site.Courses
	
	// Use a buffer to catch template execution errors before sending headers
	var buf bytes.Buffer
	if err := t.ExecuteTemplate(&buf, "layout", d); err != nil {
		log.Println("render error:", err)
		http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	buf.WriteTo(w)
}

// ---------- simple pages ----------

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	render(w, "home", http.StatusOK, pageData{Title: "Home"})
}
func AboutHandler(w http.ResponseWriter, r *http.Request) {
	render(w, "about", http.StatusOK, pageData{Title: "About"})
}
func CoursesHandler(w http.ResponseWriter, r *http.Request) {
	render(w, "courses", http.StatusOK, pageData{Title: "Courses"})
}
func SchoolHandler(w http.ResponseWriter, r *http.Request) {
	render(w, "schools", http.StatusOK, pageData{Title: "For Schools"})
}
func TestimonialHandler(w http.ResponseWriter, r *http.Request) {
	render(w, "testimonials", http.StatusOK, pageData{Title: "Testimonials"})
}

func ContactHandler(w http.ResponseWriter, r *http.Request) {
	render(w, "contact", http.StatusOK, pageData{
		Title: "Contact",
		Sent:  r.URL.Query().Get("sent") == "1",
		Form:  map[string]string{"interest": r.URL.Query().Get("interest")},
	})
}

// ---------- blog ----------

func BlogHandler(w http.ResponseWriter, r *http.Request) {
	render(w, "blog", http.StatusOK, pageData{Title: "Blog", Posts: posts})
}

func BlogPostHandler(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	for i := range posts {
		if posts[i].Slug == slug {
			render(w, "post", http.StatusOK, pageData{Title: posts[i].Title, Post: &posts[i]})
			return
		}
	}
	http.NotFound(w, r)
}

// ---------- apply / register ----------

func ApplyHandler(w http.ResponseWriter, r *http.Request) {
	render(w, "apply", http.StatusOK, pageData{
		Title: "Apply",
		Sent:  r.URL.Query().Get("sent") == "1",
		Form:  map[string]string{},
	})
}

// 0806..., 0706..., or +234806...  (Nigerian mobile)
var phoneRe = regexp.MustCompile(`^(0|\+234)[789][01]\d{8}$`)

func normalizePhone(s string) (string, bool) {
	s = strings.NewReplacer(" ", "", "-", "", "(", "", ")", "").Replace(s)
	if !phoneRe.MatchString(s) {
		return "", false
	}
	if strings.HasPrefix(s, "0") {
		s = "+234" + s[1:]
	}
	return s, true
}

func ApplySubmitHandler(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<16)
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	form := map[string]string{}
	for _, k := range []string{"name", "phone", "age", "course", "location"} {
		form[k] = strings.TrimSpace(r.PostFormValue(k))
	}

	errs := map[string]string{}
	if form["name"] == "" || len(form["name"]) > 100 {
		errs["name"] = "Enter your full name."
	}
	phone, ok := normalizePhone(form["phone"])
	if !ok {
		errs["phone"] = "Enter a valid Nigerian phone number."
	}
	age, err := strconv.Atoi(form["age"])
	if err != nil || age < 5 || age > 100 {
		errs["age"] = "Enter a valid age."
	}
	if !slices.Contains(courseTitles(), form["course"]) {
		errs["course"] = "Choose a course from the list."
	}
	if form["location"] == "" || len(form["location"]) > 100 {
		errs["location"] = "Enter your location."
	}

	if len(errs) > 0 {
		render(w, "apply", http.StatusBadRequest, pageData{Title: "Apply", Form: form, Errors: errs})
		return
	}

	app := Application{form["name"], phone, age, form["course"], form["location"], time.Now()}
	if err := saveApplication(r.Context(), app); err != nil {
		log.Println("save application:", err)
		http.Error(w, "could not save, please try again or use WhatsApp", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/apply?sent=1", http.StatusSeeOther) // 303: refresh won't resubmit
}

// ---------- contact ----------

func ContactSubmitHandler(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<16)
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	form := map[string]string{}
	for _, k := range []string{"name", "email", "phone", "interest", "message"} {
		form[k] = strings.TrimSpace(r.PostFormValue(k))
	}

	errs := map[string]string{}
	if form["name"] == "" || len(form["name"]) > 100 {
		errs["name"] = "Enter your full name."
	}
	if form["message"] == "" {
		errs["message"] = "Please enter a message."
	}
	// If they provided a phone number, make sure it's valid
	if form["phone"] != "" {
		phone, ok := normalizePhone(form["phone"])
		if !ok {
			errs["phone"] = "Enter a valid Nigerian phone number or leave blank."
		} else {
			form["phone"] = phone
		}
	}
	// If they provided an email, do a basic length check
	if form["email"] != "" && len(form["email"]) > 150 {
		errs["email"] = "Email address is too long."
	}

	if len(errs) > 0 {
		render(w, "contact", http.StatusBadRequest, pageData{Title: "Contact", Form: form, Errors: errs})
		return
	}

	msg := ContactMessage{
		Name:     form["name"],
		Email:    form["email"],
		Phone:    form["phone"],
		Interest: form["interest"],
		Message:  form["message"],
	}

	if err := saveMessage(r.Context(), msg); err != nil {
		log.Println("save message:", err)
		http.Error(w, "could not save your message, please try again or use WhatsApp", http.StatusInternalServerError)
		return
	}
	
	http.Redirect(w, r, "/contact?sent=1", http.StatusSeeOther) // 303: refresh won't resubmit
}