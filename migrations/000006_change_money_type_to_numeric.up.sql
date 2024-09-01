ALTER TABLE IF EXISTS expenses
ALTER COLUMN amount TYPE numeric USING amount::numeric;