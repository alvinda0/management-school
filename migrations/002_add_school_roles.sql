-- Add new school roles
INSERT INTO roles (name, description) VALUES 
('system', 'System administrator with full access'),
('kepala_sekolah', 'Kepala sekolah dengan akses manajemen penuh'),
('staff', 'Staff administrasi dengan akses terbatas'),
('guru', 'Guru dengan akses pembelajaran dan siswa'),
('siswa', 'Siswa dengan akses dasar')
ON CONFLICT (name) DO NOTHING;

-- Add new school-specific permissions
INSERT INTO permissions (name, description, resource, action) VALUES 
('academic.manage', 'Manage academic data', 'academic', 'manage'),
('academic.read', 'View academic data', 'academic', 'read'),
('student.manage', 'Manage student data', 'student', 'manage'),
('student.read', 'View student data', 'student', 'read'),
('teacher.manage', 'Manage teacher data', 'teacher', 'manage'),
('teacher.read', 'View teacher data', 'teacher', 'read'),
('class.manage', 'Manage class data', 'class', 'manage'),
('class.read', 'View class data', 'class', 'read'),
('report.generate', 'Generate reports', 'report', 'generate'),
('report.read', 'View reports', 'report', 'read'),
('system.manage', 'Manage system settings', 'system', 'manage')
ON CONFLICT (name) DO NOTHING;

-- Assign all permissions to system role
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM roles r, permissions p 
WHERE r.name = 'system'
ON CONFLICT DO NOTHING;

-- Assign ALL permissions to kepala sekolah role (same as system)
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM roles r, permissions p 
WHERE r.name = 'kepala_sekolah'
ON CONFLICT DO NOTHING;

-- Assign administrative permissions to staff role
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM roles r, permissions p 
WHERE r.name = 'staff' AND p.name IN (
    'user.read', 'user.update', 'academic.read',
    'student.manage', 'student.read', 'teacher.read',
    'class.read', 'report.read'
)
ON CONFLICT DO NOTHING;

-- Assign teaching permissions to guru role
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM roles r, permissions p 
WHERE r.name = 'guru' AND p.name IN (
    'user.read', 'academic.read', 'student.read',
    'class.manage', 'class.read', 'report.read'
)
ON CONFLICT DO NOTHING;

-- Assign basic permissions to siswa role
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM roles r, permissions p 
WHERE r.name = 'siswa' AND p.name IN (
    'academic.read', 'class.read', 'report.read'
)
ON CONFLICT DO NOTHING;

-- Insert default system user (password: system123)
-- Note: This is a bcrypt hash of "system123" with cost 12
INSERT INTO users (username, email, password, role_id)
SELECT 'system', 'system@gmail.com', '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewdBdXwtGtrKxQ4i', r.id
FROM roles r WHERE r.name = 'system'
ON CONFLICT (email) DO NOTHING;