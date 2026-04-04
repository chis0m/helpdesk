# Controllers

### Auth Controller

`Login` validates input, calls auth service, and sets auth cookies.  
`Signup` creates a normal user account and returns redirect target for login.  
`Refresh` validates refresh context and rotates token pair.  
`CSRFToken` returns session-bound CSRF token metadata.

### Public Auth CSRF Endpoint

`PublicAuthCSRFToken` issues a short-lived auth CSRF token from in-memory store.  
Frontend sends it in `X-CSRF-Token` when calling `/auth/login` or `/auth/signup`.

### User Controller

`UserController` contains user management endpoints.  
`UpdateRoleByUserID` is the admin role update endpoint for baseline flow.
