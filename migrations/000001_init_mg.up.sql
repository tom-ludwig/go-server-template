CREATE TABLE users (
    user_id      UUID PRIMARY KEY DEFAULT uuidv7(),
    email        TEXT,
    first_name   TEXT NOT NULL,
    last_name    TEXT NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);
