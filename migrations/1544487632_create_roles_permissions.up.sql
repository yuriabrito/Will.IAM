CREATE TABLE IF NOT EXISTS roles_permissions (
	id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
	role_id UUID,
	permission_id UUID,
	created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
	updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
  FOREIGN KEY(role_id) REFERENCES roles (id),
  FOREIGN KEY(permission_id) REFERENCES permissions (id)
);

CREATE INDEX roles_permissions_role ON roles_permissions (role_id);
CREATE INDEX roles_permissions_permission ON roles_permissions (permission_id);
