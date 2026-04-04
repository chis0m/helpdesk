// Package routes wires HTTP endpoints to controllers and middleware.
//
// |--------------------------------------------------------------------------
// | Public Auth Routes
// |--------------------------------------------------------------------------
//
// Public routes include login CSRF token issue and login endpoint.
// Login route applies login CSRF middleware.
//
// |--------------------------------------------------------------------------
// | Protected Auth Routes
// |--------------------------------------------------------------------------
//
// Refresh route applies refresh-token middleware and CSRF middleware.
// Session CSRF endpoint is behind access-token + active-session middleware.
package routes
