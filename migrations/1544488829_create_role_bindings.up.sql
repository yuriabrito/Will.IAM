CREATE TABLE IF NOT EXISTS role_bindings (
	id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
	service_account_id UUID,
	role_id UUID,
	created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
	updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
  FOREIGN KEY(service_account_id) REFERENCES service_accounts (id) ON DELETE CASCADE,
  FOREIGN KEY(role_id) REFERENCES roles (id) ON DELETE CASCADE
);

CREATE INDEX role_bindings_service_account ON role_bindings (service_account_id);

CREATE INDEX role_bindings_role ON role_bindings (role_id);
