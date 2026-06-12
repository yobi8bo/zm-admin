-- =========================================
-- 栈序管理平台初始化数据
-- 管理员账号：admin  密码：Admin@123456
-- =========================================

-- 1. 初始化部门
INSERT INTO `sys_dept` (`id`, `parent_id`, `name`, `sort`, `leader`, `phone`, `email`, `status`, `remark`, `created_at`, `updated_at`)
VALUES (1, 0, '总公司', 0, '管理员', '', '', 1, '', NOW(), NOW());

INSERT INTO `sys_dept` (`id`, `parent_id`, `name`, `sort`, `leader`, `phone`, `email`, `status`, `remark`, `created_at`, `updated_at`)
VALUES (2, 1, '技术部', 1, '', '', '', 1, '', NOW(), NOW());

INSERT INTO `sys_dept` (`id`, `parent_id`, `name`, `sort`, `leader`, `phone`, `email`, `status`, `remark`, `created_at`, `updated_at`)
VALUES (3, 1, '运营部', 2, '', '', '', 1, '', NOW(), NOW());

-- 2. 初始化角色
INSERT INTO `sys_role` (`id`, `name`, `code`, `sort`, `status`, `remark`, `created_at`, `updated_at`)
VALUES (1, '超级管理员', 'admin', 1, 1, '拥有所有权限', NOW(), NOW());

INSERT INTO `sys_role` (`id`, `name`, `code`, `sort`, `status`, `remark`, `created_at`, `updated_at`)
VALUES (2, '普通用户', 'user', 2, 1, '基础查看权限', NOW(), NOW());

-- 3. 初始化菜单（目录）
INSERT INTO `sys_menu` (`id`, `parent_id`, `name`, `type`, `path`, `component`, `permission`, `icon`, `sort`, `visible`, `status`, `created_at`, `updated_at`)
VALUES (1, 0, '首页', 2, '/dashboard', 'dashboard/index', '', 'HomeOutlined', 1, 1, 1, NOW(), NOW());

INSERT INTO `sys_menu` (`id`, `parent_id`, `name`, `type`, `path`, `component`, `permission`, `icon`, `sort`, `visible`, `status`, `created_at`, `updated_at`)
VALUES (2, 0, '系统管理', 1, '/system', 'Layout', '', 'SettingOutlined', 2, 1, 1, NOW(), NOW());

INSERT INTO `sys_menu` (`id`, `parent_id`, `name`, `type`, `path`, `component`, `permission`, `icon`, `sort`, `visible`, `status`, `created_at`, `updated_at`)
VALUES (3, 0, '日志管理', 1, '/log', 'Layout', '', 'FileTextOutlined', 3, 1, 1, NOW(), NOW());

-- 4. 初始化菜单（系统管理子菜单）
INSERT INTO `sys_menu` (`id`, `parent_id`, `name`, `type`, `path`, `component`, `permission`, `icon`, `sort`, `visible`, `status`, `created_at`, `updated_at`)
VALUES (10, 2, '用户管理', 2, 'user', 'system/user/index', 'system:user:list', 'TeamOutlined', 1, 1, 1, NOW(), NOW());

INSERT INTO `sys_menu` (`id`, `parent_id`, `name`, `type`, `path`, `component`, `permission`, `icon`, `sort`, `visible`, `status`, `created_at`, `updated_at`)
VALUES (11, 2, '角色管理', 2, 'role', 'system/role/index', 'system:role:list', 'SafetyCertificateOutlined', 2, 1, 1, NOW(), NOW());

INSERT INTO `sys_menu` (`id`, `parent_id`, `name`, `type`, `path`, `component`, `permission`, `icon`, `sort`, `visible`, `status`, `created_at`, `updated_at`)
VALUES (12, 2, '菜单管理', 2, 'menu', 'system/menu/index', 'system:menu:list', 'AppstoreOutlined', 3, 1, 1, NOW(), NOW());

INSERT INTO `sys_menu` (`id`, `parent_id`, `name`, `type`, `path`, `component`, `permission`, `icon`, `sort`, `visible`, `status`, `created_at`, `updated_at`)
VALUES (13, 2, '部门管理', 2, 'dept', 'system/dept/index', 'system:dept:list', 'ApartmentOutlined', 4, 1, 1, NOW(), NOW());

-- 5. 初始化菜单（日志管理子菜单）
INSERT INTO `sys_menu` (`id`, `parent_id`, `name`, `type`, `path`, `component`, `permission`, `icon`, `sort`, `visible`, `status`, `created_at`, `updated_at`)
VALUES (20, 3, '操作日志', 2, 'operation', 'log/operation/index', 'system:log:operation', 'ProfileOutlined', 1, 1, 1, NOW(), NOW());

INSERT INTO `sys_menu` (`id`, `parent_id`, `name`, `type`, `path`, `component`, `permission`, `icon`, `sort`, `visible`, `status`, `created_at`, `updated_at`)
VALUES (21, 3, '登录日志', 2, 'login', 'log/login/index', 'system:log:login', 'LoginOutlined', 2, 1, 1, NOW(), NOW());

-- 6. 初始化菜单（按钮权限）
-- 用户管理按钮
INSERT INTO `sys_menu` (`id`, `parent_id`, `name`, `type`, `path`, `component`, `permission`, `icon`, `sort`, `visible`, `status`, `created_at`, `updated_at`)
VALUES (100, 10, '新增用户', 3, '', '', 'system:user:add', '', 1, 1, 1, NOW(), NOW());
INSERT INTO `sys_menu` (`id`, `parent_id`, `name`, `type`, `path`, `component`, `permission`, `icon`, `sort`, `visible`, `status`, `created_at`, `updated_at`)
VALUES (101, 10, '编辑用户', 3, '', '', 'system:user:edit', '', 2, 1, 1, NOW(), NOW());
INSERT INTO `sys_menu` (`id`, `parent_id`, `name`, `type`, `path`, `component`, `permission`, `icon`, `sort`, `visible`, `status`, `created_at`, `updated_at`)
VALUES (102, 10, '删除用户', 3, '', '', 'system:user:delete', '', 3, 1, 1, NOW(), NOW());
INSERT INTO `sys_menu` (`id`, `parent_id`, `name`, `type`, `path`, `component`, `permission`, `icon`, `sort`, `visible`, `status`, `created_at`, `updated_at`)
VALUES (103, 10, '重置密码', 3, '', '', 'system:user:resetpwd', '', 4, 1, 1, NOW(), NOW());

-- 角色管理按钮
INSERT INTO `sys_menu` (`id`, `parent_id`, `name`, `type`, `path`, `component`, `permission`, `icon`, `sort`, `visible`, `status`, `created_at`, `updated_at`)
VALUES (110, 11, '新增角色', 3, '', '', 'system:role:add', '', 1, 1, 1, NOW(), NOW());
INSERT INTO `sys_menu` (`id`, `parent_id`, `name`, `type`, `path`, `component`, `permission`, `icon`, `sort`, `visible`, `status`, `created_at`, `updated_at`)
VALUES (111, 11, '编辑角色', 3, '', '', 'system:role:edit', '', 2, 1, 1, NOW(), NOW());
INSERT INTO `sys_menu` (`id`, `parent_id`, `name`, `type`, `path`, `component`, `permission`, `icon`, `sort`, `visible`, `status`, `created_at`, `updated_at`)
VALUES (112, 11, '删除角色', 3, '', '', 'system:role:delete', '', 3, 1, 1, NOW(), NOW());

-- 菜单管理按钮
INSERT INTO `sys_menu` (`id`, `parent_id`, `name`, `type`, `path`, `component`, `permission`, `icon`, `sort`, `visible`, `status`, `created_at`, `updated_at`)
VALUES (120, 12, '新增菜单', 3, '', '', 'system:menu:add', '', 1, 1, 1, NOW(), NOW());
INSERT INTO `sys_menu` (`id`, `parent_id`, `name`, `type`, `path`, `component`, `permission`, `icon`, `sort`, `visible`, `status`, `created_at`, `updated_at`)
VALUES (121, 12, '编辑菜单', 3, '', '', 'system:menu:edit', '', 2, 1, 1, NOW(), NOW());
INSERT INTO `sys_menu` (`id`, `parent_id`, `name`, `type`, `path`, `component`, `permission`, `icon`, `sort`, `visible`, `status`, `created_at`, `updated_at`)
VALUES (122, 12, '删除菜单', 3, '', '', 'system:menu:delete', '', 3, 1, 1, NOW(), NOW());

-- 部门管理按钮
INSERT INTO `sys_menu` (`id`, `parent_id`, `name`, `type`, `path`, `component`, `permission`, `icon`, `sort`, `visible`, `status`, `created_at`, `updated_at`)
VALUES (130, 13, '新增部门', 3, '', '', 'system:dept:add', '', 1, 1, 1, NOW(), NOW());
INSERT INTO `sys_menu` (`id`, `parent_id`, `name`, `type`, `path`, `component`, `permission`, `icon`, `sort`, `visible`, `status`, `created_at`, `updated_at`)
VALUES (131, 13, '编辑部门', 3, '', '', 'system:dept:edit', '', 2, 1, 1, NOW(), NOW());
INSERT INTO `sys_menu` (`id`, `parent_id`, `name`, `type`, `path`, `component`, `permission`, `icon`, `sort`, `visible`, `status`, `created_at`, `updated_at`)
VALUES (132, 13, '删除部门', 3, '', '', 'system:dept:delete', '', 3, 1, 1, NOW(), NOW());

-- 7. 初始化管理员用户（密码：Admin@123456）
INSERT INTO `sys_user` (`id`, `dept_id`, `username`, `nickname`, `password`, `avatar`, `email`, `phone`, `gender`, `status`, `remark`, `created_at`, `updated_at`)
VALUES (1, 1, 'admin', '超级管理员', '$2a$10$7P9gMycM4rNoce7nVwZIqOeSedX3Y0BesmbGw6nAbAL8vnY8UchZC', '', 'admin@zm.com', '', 1, 1, '系统初始管理员', NOW(), NOW());

-- 8. 用户绑定角色
INSERT INTO `sys_user_role` (`sys_user_id`, `sys_role_id`) VALUES (1, 1);

-- 9. 角色绑定全部菜单（超级管理员）
INSERT INTO `sys_role_menu` (`sys_role_id`, `sys_menu_id`)
SELECT 1, `id` FROM `sys_menu`;

-- 10. Casbin 权限规则
-- 用户1(admin) 绑定 admin 角色
INSERT INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES ('g', '1', 'admin', '');

-- admin 角色拥有所有接口权限（通配）
INSERT INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES ('p', 'admin', '/api/v1/*', '*');
