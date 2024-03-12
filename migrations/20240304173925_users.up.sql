CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS role
(
    id          UUID PRIMARY KEY            NOT NULL DEFAULT gen_random_uuid(),
    title       TEXT                        NOT NULL,
    slug        TEXT UNIQUE                 NOT NULL,
    description TEXT                        NOT NULL,
    created_at  TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS permission
(
    id          UUID PRIMARY KEY            NOT NULL DEFAULT gen_random_uuid(),
    title       TEXT                        NOT NULL,
    slug        TEXT UNIQUE                 NOT NULL,
    description TEXT                        NOT NULL,
    created_at  TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS role_permission
(
    role_id       UUID REFERENCES role,
    permission_id UUID REFERENCES permission,
    created_at    TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    PRIMARY KEY (role_id, permission_id)
);

CREATE TABLE IF NOT EXISTS "user"
(
    id            UUID PRIMARY KEY            NOT NULL DEFAULT gen_random_uuid(),
    name          TEXT                        NOT NULL,
    email         CITEXT UNIQUE               NOT NULL,
    password_hash BYTEA                       NOT NULL,
    activated     BOOL                        NOT NULL DEFAULT FALSE,
    role_id       UUID REFERENCES role        NOT NULL,
    created_at    TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Admin role creation
WITH admin_role_insertion AS (
    INSERT INTO role (title, slug, description)
        VALUES ('Administrator', 'admin', 'Able to manage users, content.')
        RETURNING id AS r_id),

     permissions_insertion AS (
         INSERT INTO permission (title, slug, description)
             VALUES ('View users', 'user:view:all', 'View every user details.'),
                    ('View account details', 'user:view', 'Is able to view account details.'),
                    ('Delete users', 'user:delete:all', 'Remove every user.'),
                    ('Edit user details', 'user:edit:all', 'Edit user details.')
             RETURNING id AS p_id)

INSERT
INTO role_permission (role_id, permission_id)
SELECT admin_role_insertion.r_id, permissions_insertion.p_id
FROM admin_role_insertion,
     permissions_insertion;

-- Inactive user role creation
WITH inactive_user_role_insertion AS (
    INSERT INTO role (title, slug, description)
        VALUES ('Inactive user', 'user-inactive', 'Able to activate account through activation link.')
        RETURNING id AS r_id),

     permissions_insertion AS (
         INSERT INTO permission (title, slug, description)
             VALUES ('Activate account', 'user:activate', 'Is able to active account.'),
                    ('Activate account2', 'user:activate2', 'Is able to active account.2')
             RETURNING id AS p_id)

INSERT
INTO role_permission (role_id, permission_id)
SELECT inactive_user_role_insertion.r_id, permissions_insertion.p_id
FROM inactive_user_role_insertion,
     permissions_insertion;

-- Disabled user creation
WITH disabled_user_role_insertion AS (
    INSERT INTO role (title, slug, description)
        VALUES ('Disabled user', 'user-disabled', 'Deactivated user, only able to login and view his details.')
        RETURNING id AS r_id)

INSERT
INTO role_permission (role_id, permission_id)
SELECT disabled_user_role_insertion.r_id, ((SELECT id FROM permission WHERE slug = 'user:view'))
FROM disabled_user_role_insertion