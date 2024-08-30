CREATE TABLE IF NOT EXISTS expenses (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  name text NOT NULL,
  type text NOT NULL,
  amount money NOT NULL,
  created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
  user_id uuid NOT NULL,
  version integer NOT NULL DEFAULT 1
);