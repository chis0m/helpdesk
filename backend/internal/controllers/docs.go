// Package controllers handles HTTP request/response flow.
//
// |--------------------------------------------------------------------------
// | Auth Controller
// |--------------------------------------------------------------------------
//
// Login validates input, calls auth service, and sets auth cookies.
// Refresh validates refresh context and rotates token pair.
// CSRFToken returns session-bound CSRF token metadata.
//
// |--------------------------------------------------------------------------
// | Login CSRF Endpoint
// |--------------------------------------------------------------------------
//
// LoginCSRFToken issues a short-lived login CSRF token from in-memory store.
// Frontend sends it in X-CSRF-Token when calling /auth/login.
package controllers
