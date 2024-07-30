
INSERT INTO `permissions` (`uuid`, `name`, `bit`, `parent_uuid`, `created_at`, `updated_at`)
VALUES 
(UUID(), '系统管理', 4, "", NOW(), NOW());

SET @system_uuid = (SELECT `uuid` FROM `permissions` WHERE `name` = '系统管理');

INSERT INTO `permissions` (`uuid`, `name`, `bit`, `parent_uuid`, `created_at`, `updated_at`)
VALUES
(UUID(), '用户管理', 4, @system_uuid, NOW(), NOW()),
(UUID(), '角色管理', 4, @system_uuid, NOW(), NOW()),
(UUID(), '菜单管理', 4, @system_uuid, NOW(), NOW()),
(UUID(), '登录日志', 4, @system_uuid, NOW(), NOW()),
(UUID(), '操作日志', 4, @system_uuid, NOW(), NOW()),
(UUID(), 'API管理', 4, @system_uuid, NOW(), NOW()),
(UUID(), '权限管理', 4, @system_uuid, NOW(), NOW());


SET @user_uuid = (SELECT `uuid` FROM `permissions` WHERE `name` = '用户管理');
SET @role_uuid = (SELECT `uuid` FROM `permissions` WHERE `name` = '角色管理');
SET @menu_uuid = (SELECT `uuid` FROM `permissions` WHERE `name` = '菜单管理');
SET @loginlog_uuid = (SELECT `uuid` FROM `permissions` WHERE `name` = '登录日志');
SET @oplog_uuid = (SELECT `uuid` FROM `permissions` WHERE `name` = '操作日志');
SET @api_uuid = (SELECT `uuid` FROM `permissions` WHERE `name` = 'API管理');
SET @permission_uuid = (SELECT `uuid` FROM `permissions` WHERE `name` = '权限管理');


INSERT INTO `permissions` (`uuid`, `name`, `bit`, `parent_uuid`, `created_at`, `updated_at`)
VALUES
(UUID(), '创建用户', 8, @user_uuid, NOW(), NOW()),
(UUID(), '更新用户信息', 2, @user_uuid, NOW(), NOW()),
(UUID(), '删除用户', 1, @user_uuid, NOW(), NOW());


INSERT INTO `permissions` (`uuid`, `name`, `bit`, `parent_uuid`, `created_at`, `updated_at`)
VALUES
(UUID(), '创建角色', 8, @role_uuid, NOW(), NOW()),
(UUID(), '编辑角色', 2, @role_uuid, NOW(), NOW()),
(UUID(), '删除角色', 1, @role_uuid, NOW(), NOW());


INSERT INTO `permissions` (`uuid`, `name`, `bit`, `parent_uuid`, `created_at`, `updated_at`)
VALUES
(UUID(), '创建菜单', 8, @menu_uuid, NOW(), NOW()),
(UUID(), '编辑菜单', 2, @menu_uuid, NOW(), NOW()),
(UUID(), '删除菜单', 1, @menu_uuid, NOW(), NOW());


INSERT INTO `permissions` (`uuid`, `name`, `bit`, `parent_uuid`, `created_at`, `updated_at`)
VALUES
(UUID(), '创建API', 8, @api_uuid, NOW(), NOW()),
(UUID(), '编辑API', 2, @api_uuid, NOW(), NOW()),
(UUID(), '删除API', 1, @api_uuid, NOW(), NOW());


INSERT INTO `permissions` (`uuid`, `name`, `bit`, `parent_uuid`, `created_at`, `updated_at`)
VALUES
(UUID(), '创建权限', 8, @permission_uuid, NOW(), NOW()),
(UUID(), '编辑权限', 2, @permission_uuid, NOW(), NOW()),
(UUID(), '删除权限', 1, @permission_uuid, NOW(), NOW());