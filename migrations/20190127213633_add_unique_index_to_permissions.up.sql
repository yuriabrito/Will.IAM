CREATE UNIQUE INDEX IF NOT EXISTS permissions_unique ON permissions (role_id, ownership_level, action, service, resource_hierarchy);
