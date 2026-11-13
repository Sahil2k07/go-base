CREATE UNIQUE INDEX IF NOT EXISTS uq_roles_name 
ON roles (name) 
WHERE deleted_at IS NULL;