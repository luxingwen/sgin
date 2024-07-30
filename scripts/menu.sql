-- 插入一级菜单
INSERT INTO `menus` (`uuid`, `name`, `link`, `parent_uuid`, `created_at`, `updated_at`, `icon`, `order`, `is_show`, `type`)
VALUES 
(UUID(), '首页', '/home', NULL, NOW(), NOW(), 'home', 1, 1, 1),
(UUID(), '登录', '/user/login', NULL, NOW(), NOW(), 'login', 2, 0, 2),
(UUID(), '个人中心', '/user/profile', NULL, NOW(), NOW(), 'user', 3, 0, 2),
(UUID(), '系统管理', '/system', NULL, NOW(), NOW(), 'setting', 9, 1, 1);


SET @system_uuid = (SELECT `uuid` FROM `menus` WHERE `name` = '系统管理');

-- 插入二级菜单
INSERT INTO `menus` (`uuid`, `name`, `link`, `parent_uuid`, `created_at`, `updated_at`, `icon`, `order`, `is_show`, `type`)
VALUES
(UUID(), '用户管理', '/system/user', @system_uuid, NOW(), NOW(), 'user', 1, 1, 1),
(UUID(), '角色管理', '/system/role', @system_uuid, NOW(), NOW(), 'team', 2, 1, 1),
(UUID(), '菜单管理', '/system/menu', @system_uuid, NOW(), NOW(), 'menu', 3, 1, 1),
(UUID(), '登录日志', '/system/loginlog', @system_uuid, NOW(), NOW(), 'login', 4, 1, 1),
(UUID(), '操作日志', '/system/oplog', @system_uuid, NOW(), NOW(), 'audit', 5, 1, 1),
(UUID(), 'API管理', '/system/api', @system_uuid, NOW(), NOW(), 'api', 6, 1, 1),
(UUID(), '权限管理', '/system/permission', @system_uuid, NOW(), NOW(), '', 7, 1, 1);
