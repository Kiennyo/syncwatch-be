DELETE FROM role_permission WHERE role_id = (SELECT id FROM role WHERE slug = 'user-active');
DELETE FROM permission WHERE slug IN ('user:edit');
DELETE FROM role WHERE id = (SELECT id FROM role WHERE slug = 'user-active');
