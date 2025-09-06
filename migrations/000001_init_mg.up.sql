CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
  user_id      UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  email        TEXT,
  first_name   TEXT NOT NULL,
  last_name    TEXT NOT NULL
);
