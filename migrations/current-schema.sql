-- Use this to compare to a old schema state that you think is the current database state
CREATE TABLE testusers (
    user_id      UUID PRIMARY KEY DEFAULT uuidv7(),
    email        TEXT,
    first_name   TEXT NOT NULL,
    last_name    TEXT NOT NULL,
    thrind_name    TEXT NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);
