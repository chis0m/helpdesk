// Package services contains business rules between controllers and repositories.
//
// |--------------------------------------------------------------------------
// | Auth Service
// |--------------------------------------------------------------------------
//
// AuthService handles:
// - login credential verification
// - access/refresh token creation
// - refresh token rotation and replay checks
// - session-bound CSRF token issue/reuse
//
// |--------------------------------------------------------------------------
// | User Service
// |--------------------------------------------------------------------------
//
// UserService handles create/update user operations.
package services
