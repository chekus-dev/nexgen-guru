CREATE TABLE IF NOT EXISTS page_views (
    id         BIGSERIAL PRIMARY KEY,
    path       TEXT        NOT NULL, -- e.g., "/courses", "/apply"
    ip_address TEXT        NOT NULL DEFAULT '',
    user_agent TEXT        NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);