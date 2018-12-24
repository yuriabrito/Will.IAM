CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS service_accounts (
	id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
  name VARCHAR(200) NOT NULL,
	key_id VARCHAR(200),
	key_secret VARCHAR(200),
	email VARCHAR(200),
  base_role_id UUID NOT NULL,
	created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
	updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX IF NOT EXISTS service_accounts_name ON service_accounts (name);
CREATE UNIQUE INDEX IF NOT EXISTS service_accounts_key_id_key_secret ON service_accounts (key_id, key_secret);
CREATE UNIQUE INDEX IF NOT EXISTS service_accounts_email ON service_accounts (email);
