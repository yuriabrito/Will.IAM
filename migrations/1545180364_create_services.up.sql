CREATE TABLE IF NOT EXISTS services (
	id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
	name VARCHAR(100) NOT NULL,
	permission_name VARCHAR(100) NOT NULL,
  service_account_id UUID NOT NULL,
  creator_service_account_id UUID NOT NULL,
	created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
	updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
  FOREIGN KEY(service_account_id) REFERENCES service_accounts (id) ON DELETE CASCADE
);
