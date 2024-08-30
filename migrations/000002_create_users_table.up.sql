CREATE TABLE IF NOT EXISTS users (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at timestamp(0) with time zone DEFAULT NOW(),
  name text NOT NULL,
  email text UNIQUE NOT NULL,
  password_hash bytea NOT NULL,
  version integer NOT NULL DEFAULT 1
);