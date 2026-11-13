CREATE UNIQUE INDEX IF NOT EXISTS uq_users_email
ON users (email) 
WHERE deleted_at IS NULL;