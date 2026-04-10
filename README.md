# Helpdesk App

This project is a secure web helpdesk application for managing support users and requests.

It has:

- `frontend` (Vue + TypeScript UI)
- `backend` (Go API + database logic)

Database: `MySQL`.

Backend architecture: `MVC` with dependency injection.

Password hashing: `Argon2id`.

Logging: `zerolog` (pretty logs in development, structured logs in production).

## Documentation conventions

Documentation is colocated as `*.md` files next to the main backend and frontend areas. 
The outline is consistent: what that area covers, then how requests and data move through it, 
so the docs read the same wherever you open them.