-- Add role directly to users
ALTER TABLE users ADD COLUMN role TEXT NOT NULL DEFAULT 'user';

ALTER TABLE users ADD CONSTRAINT role_check CHECK (role IN ('user', 'admin', 'moderator'));

-- Permissions table
CREATE TABLE permissions (
    id SERIAL PRIMARY KEY,
    codename TEXT UNIQUE NOT NULL,
    name TEXT NOT NULL
);

-- Groups table
CREATE TABLE groups (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL
);

-- Group-Permissions (Many-to-Many)
CREATE TABLE group_permissions (
    group_id INTEGER NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
    permission_id INTEGER NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    PRIMARY KEY (group_id, permission_id)
);

-- User-Groups (Many-to-Many)
CREATE TABLE user_groups (
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    group_id INTEGER NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, group_id)
);
