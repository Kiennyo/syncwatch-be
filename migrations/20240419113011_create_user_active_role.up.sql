WITH user_active_role AS (
    INSERT INTO role (title, slug, description)
        VALUES ('Activated user', 'user-active', 'Able to view profile, edit, upload videos and create sharable links')
        RETURNING id AS r_id),

     permissions_insertion AS (
         INSERT INTO permission (title, slug, description)
             VALUES ('Edit account details', 'user:edit', 'Is able to edit account details')
             RETURNING id AS p_id)

INSERT
INTO role_permission (role_id, permission_id)
SELECT user_active_role.r_id, permissions_insertion.p_id
FROM user_active_role,
     permissions_insertion
UNION
SELECT user_active_role.r_id, (SELECT id FROM permission WHERE slug IN ('user:view'))
FROM user_active_role;