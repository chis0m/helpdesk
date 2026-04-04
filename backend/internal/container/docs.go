// Package container wires dependencies used by the application.
//
// |--------------------------------------------------------------------------
// | Dependency Wiring
// |--------------------------------------------------------------------------
//
// Container builds and holds shared instances:
// repositories, services, controllers, token maker, and CSRF store.
//
// |--------------------------------------------------------------------------
// | Login CSRF Store
// |--------------------------------------------------------------------------
//
// LoginCSRFStore is initialized in memory for development/CA usage.
// Replace with Redis/shared storage in production.
package container
