CREATE TABLE testusers (
    user_id      UUID PRIMARY KEY DEFAULT uuidv7(),
    email        TEXT,
    first_name   TEXT,
    last_name    TEXT, 
    thrind_name  TEXT, 
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);
