INSERT INTO `apis` (`uuid`, `module`, `name`, `permission_level`, `status`, `path`, `method`, `created_at`, `updated_at`) VALUES 
(UUID(), '用户信息', '创建用户', 2, 1, '/api/v1/user/create', 'POST', NOW(), NOW()),
(UUID(), '用户信息', '获取用户信息', 2, 1, '/api/v1/user/info', 'POST', NOW(), NOW()),
(UUID(), '用户信息', '获取用户列表', 2, 1, '/api/v1/user/list', 'POST', NOW(), NOW()),
(UUID(), '用户信息', '更新用户信息', 2, 1, '/api/v1/user/update', 'POST', NOW(), NOW()),
(UUID(), '用户信息', '删除用户', 2, 1, '/api/v1/user/delete', 'POST', NOW(), NOW()),
(UUID(), '用户信息', '获取我的信息', 2, 1, '/api/v1/user/myinfo', 'GET', NOW(), NOW()),
(UUID(), '用户信息', '上传头像', 2, 1, '/api/v1/user/avatar', 'POST', NOW(), NOW());


INSERT INTO `apis` (`uuid`, `module`, `name`, `permission_level`, `status`, `path`, `method`, `created_at`, `updated_at`) VALUES 
(UUID(), '菜单', '创建菜单', 2, 1, '/api/v1/menu/create', 'POST', NOW(), NOW()),
(UUID(), '菜单', '获取菜单列表', 2, 1, '/api/v1/menu/list', 'POST', NOW(), NOW()),
(UUID(), '菜单', '更新菜单', 2, 1, '/api/v1/menu/update', 'POST', NOW(), NOW()),
(UUID(), '菜单', '删除菜单', 2, 1, '/api/v1/menu/delete', 'POST', NOW(), NOW());


INSERT INTO `apis` (`uuid`, `module`, `name`, `permission_level`, `status`, `path`, `method`, `created_at`, `updated_at`) VALUES 
(UUID(), '通用', '用户登录', 1, 1, '/api/v1/login', 'POST', NOW(), NOW());
