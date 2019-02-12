DROP INDEX IF EXISTS tokens_sso_access_token;
DROP INDEX IF EXISTS tokens_refresh_token;
ALTER TABLE tokens DROP COLUMN sso_access_token;
ALTER TABLE tokens ADD COLUMN expired_at TIMESTAMP;
