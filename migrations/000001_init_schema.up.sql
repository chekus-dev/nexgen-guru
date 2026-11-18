-- Applications Table (from /apply)
CREATE TABLE IF NOT EXISTS applications (
    id         BIGSERIAL PRIMARY KEY,
    name       TEXT        NOT NULL,
    phone      TEXT        NOT NULL,
    age        INT         NOT NULL CHECK (age BETWEEN 5 AND 100),
    course     TEXT        NOT NULL,
    location   TEXT        NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS applications_created_idx ON applications (created_at DESC);
ALTER TABLE applications ENABLE ROW LEVEL SECURITY;

-- Contact Messages Table (from /contact)
CREATE TABLE IF NOT EXISTS messages (
    id         BIGSERIAL PRIMARY KEY,
    name       TEXT        NOT NULL,
    email      TEXT        NOT NULL DEFAULT '',
    phone      TEXT        NOT NULL DEFAULT '',
    interest   TEXT        NOT NULL DEFAULT '',
    message    TEXT        NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

ALTER TABLE messages ENABLE ROW LEVEL SECURITY;