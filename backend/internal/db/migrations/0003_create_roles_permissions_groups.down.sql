-- Drop group_permissions table
DROP TABLE IF EXISTS group_permissions;

-- Drop user_groups table
DROP TABLE IF EXISTS user_groups;

-- Drop groups table
DROP TABLE IF EXISTS groups;

-- Drop permissions table
DROP TABLE IF EXISTS permissions;

-- Drop roles enum from users table and remove the column
ALTER TABLE users DROP CONSTRAINT IF EXISTS role_check;

ALTER TABLE users DROP COLUMN role;

DROP TYPE IF EXISTS user_role;
