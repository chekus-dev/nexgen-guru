CREATE TABLE IF NOT EXISTS subscribers (
    id         BIGSERIAL PRIMARY KEY,
    email      TEXT        NOT NULL UNIQUE, -- Prevents duplicate emails
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

ALTER TABLE subscribers ENABLE ROW LEVEL SECURITY;