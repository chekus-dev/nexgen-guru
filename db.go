package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var db *sql.DB

func connectDB() error {
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		return errors.New("DATABASE_URL is not set in .env file")
	}

	var err error
	db, err = sql.Open("pgx", url)
	if err != nil {
		return err
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	log.Println("Connecting to database...")
	if err := db.PingContext(ctx); err != nil {
		return err
	}
	log.Println("✅ Successfully connected to the database!")
	return nil
}

// --- APPLICATIONS ---
func saveApplication(ctx context.Context, a Application) error {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	_, err := db.ExecContext(ctx,
		`INSERT INTO applications (name, phone, age, course, location) 
		 VALUES ($1, $2, $3, $4, $5)`,
		a.Name, a.Phone, a.Age, a.Course, a.Location)
	
	if err != nil {
		log.Printf("SQL Error (applications): %v\n", err)
	}
	return err
}

// --- CONTACT MESSAGES ---
type ContactMessage struct{ 
	Name, Email, Phone, Interest, Message string 
}

func saveMessage(ctx context.Context, m ContactMessage) error {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	_, err := db.ExecContext(ctx,
		`INSERT INTO messages (name, email, phone, interest, message) 
		 VALUES ($1, $2, $3, $4, $5)`,
		m.Name, m.Email, m.Phone, m.Interest, m.Message)
	
	if err != nil {
		log.Printf("SQL Error (messages): %v\n", err)
	}
	return err
}

// --- NEWSLETTER SUBSCRIBERS ---
func saveSubscriber(ctx context.Context, email string) error {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	// The "ON CONFLICT DO NOTHING" prevents errors if the email already exists.
	_, err := db.ExecContext(ctx,
		`INSERT INTO subscribers (email) VALUES ($1) ON CONFLICT (email) DO NOTHING`,
		email)
	
	if err != nil {
		log.Printf("SQL Error (subscribers): %v\n", err)
	}
	return err
}


func saveTestimonial(ctx context.Context, t Testimonial) error {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	_, err := db.ExecContext(ctx,
		`INSERT INTO testimonials (name, course, quote, rating) 
		 VALUES ($1, $2, $3, $4)`,
		t.Name, t.Course, t.Quote, t.Rating)
	
	if err != nil {
		log.Printf("SQL Error (testimonials): %v\n", err)
	}
	return err
}

// --- FETCH TESTIMONIALS (NEW) ---
// getTestimonials retrieves all public testimonials from the database,
// ordered by the most recent first.
func getTestimonials(ctx context.Context) ([]Testimonial, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rows, err := db.QueryContext(ctx,
		`SELECT id, name, course, quote, rating 
		 FROM testimonials 
		 WHERE is_public = true 
		 ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var testimonials []Testimonial
	for rows.Next() {
		var t Testimonial
		if err := rows.Scan(&t.ID, &t.Name, &t.Course, &t.Quote, &t.Rating); err != nil {
			return nil, err
		}
		testimonials = append(testimonials, t)
	}
	
	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		return nil, err
	}
	
	return testimonials, nil
}