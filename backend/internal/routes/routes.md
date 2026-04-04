# Routes

### Public Auth Routes

Public routes include auth CSRF token issue, login, and signup.  
Login/signup routes apply public-auth CSRF middleware and rate limits.

### Protected Auth Routes

Refresh route applies refresh-token middleware and CSRF middleware.  
Session CSRF endpoint is behind access-token + active-session middleware.  
Admin role update route currently uses `user_id` path parameter in baseline.
