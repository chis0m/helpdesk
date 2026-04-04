# Container

### Dependency Wiring

Container builds and holds shared instances: repositories, services, controllers, token maker, and public auth CSRF store.

### How Container Works

`container.New(db, cfg, tokenMaker)` is called during app boot.

Inside `New`:

1. Repositories are created first (for database access):
   - `UserRepository`
   - `AuthSessionRepository`
2. Shared in-memory stores are created:
   - `PublicAuthCSRFStore`
3. Services are created next and receive repositories:
   - `UserService`
   - `AuthService`
4. Controllers are created last and receive services/stores:
   - `HealthController`
   - `AuthController`
   - `UserController`
5. All instances are stored in `Container` struct fields.

### Runtime Usage

Routes do not build dependencies directly.  
Routes receive a single `Container` and pull what they need:

- controller handlers (`AuthController`, `UserController`, `HealthController`)
- middleware dependencies (`TokenMaker`, `SessionRepo`, `PublicAuthCSRFStore`)

This keeps wiring in one place and avoids creating services/controllers inside route files.

### Public Auth CSRF Store

`PublicAuthCSRFStore` is initialized in memory for development/CA usage.  
For production, I use Redis/shared storage.
