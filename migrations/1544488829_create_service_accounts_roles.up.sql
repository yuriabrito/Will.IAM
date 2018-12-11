CREATE TABLE IF NOT EXISTS service_accounts_roles (
	id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
	service_account_id UUID,
	role_id UUID,
	created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
	updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
  FOREIGN KEY(service_account_id) REFERENCES service_accounts (id),
  FOREIGN KEY(role_id) REFERENCES roles (id)
);

CREATE INDEX service_accounts_roles_service_account ON service_accounts_roles (service_account_id);

CREATE INDEX service_accounts_roles_role ON service_accounts_roles (role_id);
