CREATE UNIQUE INDEX IF NOT EXISTS uq_permissions_name 
ON permissions (name) 
WHERE deleted_at IS NULL;