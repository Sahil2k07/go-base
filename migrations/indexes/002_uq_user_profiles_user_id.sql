CREATE UNIQUE INDEX IF NOT EXISTS uq_user_profiles_user_id 
ON user_profiles (user_id) 
WHERE deleted_at IS NULL;