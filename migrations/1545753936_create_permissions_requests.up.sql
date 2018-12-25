CREATE TABLE IF NOT EXISTS permissions_requests (
	id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
	service VARCHAR(200) NOT NULL,
	action VARCHAR(200) NOT NULL,
	resource_hierarchy VARCHAR(200) NOT NULL,
  message VARCHAR(200) NOT NULL,
  state SMALLINT NOT NULL,
	service_account_id UUID NOT NULL,
	created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
	updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
  FOREIGN KEY(service_account_id) REFERENCES service_accounts (id) ON DELETE CASCADE
);

CREATE INDEX permissions_requests_service_account ON permissions_requests (service_account_id);
