-- Migration to change from ID to UUID and remove permissions system

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Drop existing foreign key constraints
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_role_id_fkey;
ALTER TABLE role_permissions DROP CONSTRAINT IF EXISTS role_permissions_role_id_fkey;
ALTER TABLE role_permissions DROP CONSTRAINT IF EXISTS role_permissions_permission_id_fkey;

-- Drop permissions-related tables
DROP TABLE IF EXISTS role_permissions;
DROP TABLE IF EXISTS permissions;

-- Create new tables with UUID
CREATE TABLE IF NOT EXISTS roles_new (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS users_new (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    role_id UUID NOT NULL REFERENCES roles_new(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

-- Migrate data from old tables to new tables
INSERT INTO roles_new (name, description, created_at, updated_at)
SELECT name, description, created_at, updated_at FROM roles;

-- For users, we need to map the old role_id to new UUID role_id
INSERT INTO users_new (username, email, password, role_id, created_at, updated_at, deleted_at)
SELECT u.username, u.email, u.password, rn.id, u.created_at, u.updated_at, u.deleted_at
FROM users u
JOIN roles r ON u.role_id = r.id
JOIN roles_new rn ON r.name = rn.name;

-- Drop old tables
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS roles;

-- Rename new tables
ALTER TABLE roles_new RENAME TO roles;
ALTER TABLE users_new RENAME TO users;

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);
CREATE INDEX IF NOT EXISTS idx_users_role_id ON users(role_id);