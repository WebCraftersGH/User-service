CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  username VARCHAR(25),
  email VARCHAR(100) NOT NULL,
  fio VARCHAR(100),
  bio VARCHAR(255),
  sex TEXT CHECK (sex IN ('Male', 'Female', 'Other')),
  birthday DATE,
  last_login_date timestamptz,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  deleted_at timestamptz
);

CREATE UNIQUE INDEX IF NOT EXISTS ux_users_email ON users(email) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS ix_users_deleted_at ON users(deleted_at);
