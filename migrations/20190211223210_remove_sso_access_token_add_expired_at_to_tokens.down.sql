ALTER TABLE tokens ADD COLUMN sso_access_token VARCHAR(300) NOT NULL;
ALTER TABLE tokens DROP COLUMN expired_at;
CREATE UNIQUE INDEX tokens_sso_access_token ON tokens (sso_access_token);
CREATE UNIQUE INDEX tokens_refresh_token ON tokens (refresh_token);
