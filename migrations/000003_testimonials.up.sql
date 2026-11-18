CREATE TABLE IF NOT EXISTS testimonials (
    id         BIGSERIAL PRIMARY KEY,
    name       TEXT        NOT NULL,
    course     TEXT        NOT NULL,
    quote      TEXT        NOT NULL,
    rating     INT         NOT NULL DEFAULT 5 CHECK (rating BETWEEN 1 AND 5),
    is_public  BOOLEAN     NOT NULL DEFAULT true, -- Allows you to hide a testimonial without deleting it
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

ALTER TABLE testimonials ENABLE ROW LEVEL SECURITY;