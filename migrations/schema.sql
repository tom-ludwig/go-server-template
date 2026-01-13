-- Schema definition for go-server-template
-- This file is used with pg-schema-diff for declarative migrations

CREATE TABLE users (
    user_id      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email        TEXT,
    first_name   TEXT NOT NULL,
    last_name    TEXT NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);
