ALTER TABLE IF EXISTS expenses
ALTER COLUMN amount TYPE money USING amount::money;