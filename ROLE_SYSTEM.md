# Sistem Role dan Permission

## Overview
Sistem ini menggunakan Role-Based Access Control (RBAC) dengan permission yang fleksibel. Role "system" dan "kepala_sekolah" memiliki akses penuh ke semua fitur.

## Role Hierarchy

### 1. System
- **Akses**: SEMUA permission dan fitur
- **Deskripsi**: Administrator sistem dengan kontrol penuh
- **Default User**: system@gmail.com (password: system123)

### 2. Kepala Sekolah
- **Akses**: SEMUA permission dan fitur (sama seperti system)
- **Deskripsi**: Kepala sekolah dengan akses manajemen penuh
- **Dapat mengakses**: Semua endpoint dan fitur tanpa batasan

### 3. Staff
- **Akses**: Administrative permissions
- **Permissions**:
  - user.read, user.update
  - academic.read
  - student.manage, student.read
  - teacher.read
  - class.read
  - report.read

### 4. Guru
- **Akses**: Teaching permissions
- **Permissions**:
  - user.read
  - academic.read
  - student.read
  - class.manage, class.read
  - report.read

### 5. Siswa
- **Akses**: Basic permissions
- **Permissions**:
  - academic.read
  - class.read
  - report.read

## Middleware Functions

### 1. AuthMiddleware()
- Validasi JWT token
- Set userID dan roleID di context

### 2. RequirePermission(permission string)
- Check permission spesifik
- **System dan Kepala Sekolah**: Otomatis dapat akses SEMUA permission
- Role lain: Check permission dari database

### 3. RequireRole(roleName string)
- Check role spesifik
- **System dan Kepala Sekolah**: Dapat akses semua role-based endpoint
- Role lain: Harus sesuai dengan role yang diminta

### 4. RequireAnyRole(roleNames ...string)
- Check apakah user memiliki salah satu dari role yang diminta
- **System dan Kepala Sekolah**: Otomatis dapat akses
- Role lain: Check apakah role ada dalam daftar yang diizinkan

### 5. IsSystemOrKepalaSekolah()
- Khusus untuk endpoint yang hanya boleh diakses system dan kepala sekolah
- Endpoint admin-only

## Contoh Penggunaan

### Admin-Only Routes (System & Kepala Sekolah)
```go
adminOnly := protected.Group("/admin")
adminOnly.Use(middleware.IsSystemOrKepalaSekolah())
{
    adminOnly.POST("/users", userHandler.Register)
    adminOnly.POST("/roles", roleHandler.CreateRole)
    adminOnly.PUT("/roles/:id", roleHandler.UpdateRole)
    adminOnly.DELETE("/roles/:id", roleHandler.DeleteRole)
}
```

### Management Routes (Multiple Roles)
```go
management := protected.Group("/management")
management.Use(middleware.RequireAnyRole("kepala_sekolah", "staff", "guru"))
{
    management.GET("/academic", middleware.RequirePermission("academic.read"), handler)
    management.POST("/students", middleware.RequirePermission("student.manage"), handler)
}
```

### Permission-Based Routes
```go
protected.GET("/users", middleware.RequirePermission("user.read"), userHandler.GetAllUsers)
protected.PUT("/users/:id", middleware.RequirePermission("user.update"), userHandler.UpdateUser)
```

## Keunggulan Sistem

1. **Fleksibilitas**: System dan Kepala Sekolah dapat akses semuanya tanpa perlu update permission
2. **Granular Control**: Permission level yang detail untuk role lain
3. **Scalable**: Mudah menambah role dan permission baru
4. **Database-Driven**: Permission check menggunakan database, bukan hardcoded
5. **Hierarchical**: Role dengan level akses yang jelas

## Testing

### Login sebagai System Admin
```bash
POST /api/v1/login
{
    "email": "system@gmail.com",
    "password": "system123"
}
```

### Test Access
Dengan token system atau kepala sekolah, semua endpoint akan dapat diakses:
- GET /api/v1/users ✅
- POST /api/v1/admin/users ✅
- DELETE /api/v1/admin/roles/1 ✅
- POST /api/v1/management/students ✅
- GET /api/v1/teacher/classes ✅ (meskipun bukan guru)

## Migration & Seeding

Jalankan migration dan seeding untuk setup role dan permission:
```bash
go run cmd/migrate/main.go
go run cmd/seed/main.go
```

Atau gunakan setup command:
```bash
go run cmd/setup/main.go
```