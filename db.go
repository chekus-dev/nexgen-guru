package main

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib" // registers the "pgx" driver
)

var db *sql.DB

// Run in order on every startup. IF NOT EXISTS makes them safe to repeat.
var schema = []string{
	`CREATE TABLE IF NOT EXISTS applications (
		id         BIGSERIAL PRIMARY KEY,
		name       TEXT        NOT NULL,
		phone      TEXT        NOT NULL,
		age        INT         NOT NULL CHECK (age BETWEEN 5 AND 100),
		course     TEXT        NOT NULL,
		location   TEXT        NOT NULL,
		created_at TIMESTAMPTZ NOT NULL DEFAULT now()
	)`,
	`CREATE INDEX IF NOT EXISTS applications_created_idx ON applications (created_at DESC)`,
	// Supabase exposes tables through a public API; this keeps them private. The app's own connection is unaffected.
	`ALTER TABLE applications ENABLE ROW LEVEL SECURITY`,
	`CREATE TABLE IF NOT EXISTS messages (
		id         BIGSERIAL PRIMARY KEY,
		name       TEXT        NOT NULL,
		email      TEXT        NOT NULL DEFAULT '',
		phone      TEXT        NOT NULL DEFAULT '',
		interest   TEXT        NOT NULL DEFAULT '',
		message    TEXT        NOT NULL,
		created_at TIMESTAMPTZ NOT NULL DEFAULT now()
	)`,
	`ALTER TABLE messages ENABLE ROW LEVEL SECURITY`,
}

func connectDB() error {
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		return errors.New("DATABASE_URL is not set")
	}
	var err error
	db, err = sql.Open("pgx", url)
	if err != nil {
		return err
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return err
	}
	for _, q := range schema {
		if _, err := db.ExecContext(ctx, q); err != nil {
			return err
		}
	}
	return nil
}

func saveApplication(ctx context.Context, a Application) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	_, err := db.ExecContext(ctx,
		`INSERT INTO applications (name, phone, age, course, location) VALUES ($1, $2, $3, $4, $5)`,
		a.Name, a.Phone, a.Age, a.Course, a.Location)
	return err
}

type ContactMessage struct{ Name, Email, Phone, Interest, Message string }

func saveMessage(ctx context.Context, m ContactMessage) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	_, err := db.ExecContext(ctx,
		`INSERT INTO messages (name, email, phone, interest, message) VALUES ($1, $2, $3, $4, $5)`,
		m.Name, m.Email, m.Phone, m.Interest, m.Message)
	return err
}
