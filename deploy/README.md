# Deployment overview

This folder describes how the app runs in production: **Docker** on an **EC2** instance, with **GitHub Actions** building images and rolling out updates when you push to specific branches.

---

## What triggers a deploy

Pushes to:

| Branch | What gets deployed |
|--------|---------------------|
| `vulnerable-baseline` | **Vuln** stack  |
| `secure-fix` | **Secure** stack  |

The deployment is triggered based on the branch you push or merge to.

---

## What the pipeline does (high level)

1. **Backend configuration** — The workflow loads environment variables from **AWS Secrets Manager** (different secret per stack). Those values are used when building the backend image so settings like database credentials are baked in for that variant. Keep the secret in sync with the committed `.env.*.local` files (or your intended production values); the workflow **overwrites** those files in CI before each build.

2. **Build & push images** — Backend and frontend **Docker** images are built in CI, tagged `vuln` or `secure`, plus the git commit SHA, and pushed to **Amazon ECR Public** so the server can pull them by name.

3. **Frontend API URL** — The frontend is built with a fixed **HTTPS** API URL for that stack.

4. **Deploy to EC2 over SSH** — The workflow connects to the instance using SSH private key and server IP in  **repository secret**. It copies the **Compose** files from this repo and runs `docker compose pull` and `up` so the running containers match the new images.

---

## Domains and URLs (as configured in this repo)

The workflow and Compose files assume a fixed production domain (`chisomejim.site`) with subdomains for:

| Resource | URL |
|----------|-----|
| Vuln (UI) | https://vulnweb.chisomejim.site |
| Vuln (API) | https://api.vulnweb.chisomejim.site |
| Secure (UI) | https://secweb.chisomejim.site |
| Secure (API) | https://api.secweb.chisomejim.site |
| Mail (Mailpit UI) | https://mail.chisomejim.site |

DNS for those names point to the EC2 instance’s public IP. If domain or subdomains is changed, then workflow is updated.

---

## Layout on the server

Typical paths under the `ubuntu` user’s home directory:

| Path | Purpose |
|------|---------|
| `~/helpdesk/mail/` | **Mailpit** dev smtp server. Persists data under `mail/data/`. |
| `~/helpdesk/database/` | **MySQL 8** (Docker) for both stacks: small footprint (`innodb_buffer_pool_size` 128M, 512M container cap), hostname **`mysql`** on `helpdesk-shared`. Persist data under `database/data/`. |
| `~/helpdesk/vuln/` or `~/helpdesk/secure/` | The **app stack** for that variant (backend + web containers). |

The **mail** stack is updated on every deploy, then the **database** stack, then the **app** stack for the branch you pushed.

---

## Docker networking

- The **mail** Compose file defines a shared Docker network called **`helpdesk-shared`**.
- Each **app** stack joins that network as **external** so backends can reach Mailpit by hostname (`mailpit`) on the internal SMTP port and MySQL by hostname **`mysql`** on port **3306**.
- App HTTP ports are bound to **`127.0.0.1` only** on the host. They are not meant to be reached directly from the internet; traffic is expected to come through a **reverse proxy** on the same machine (see below).
- The database Compose file publishes **MySQL on host port 3306** (all interfaces) so clients outside the host can connect if the **EC2 security group** allows TCP 3306 from your IP. That is convenient for tooling but increases exposure; tighten the security group and use non-default passwords outside class demos.

---

## Reverse proxy (Nginx)

**Nginx** terminates HTTPS and proxies to the Docker services on `127.0.0.1` (see `nginx/chisomejim.site.conf` for hostname → port mapping). Containers are not exposed directly on the public internet.

### TLS (SSL)

HTTPS uses **Certbot** with **Let’s Encrypt** (e.g. `certbot --nginx`). Certificates and keys live on the instance under **`/etc/letsencrypt/`**; Nginx `ssl_certificate` / `ssl_certificate_key` paths must matches what **`certbot certificates`** shows.

---

## GitHub configuration

The workflow expects:

- **Secrets:** AWS credentials (access key, secret, region), SSH private key for deploy, and **`SERVER_IP`** — the EC2 instance address used for SSH.

Backend secrets live in **AWS Secrets Manager** under the IDs referenced in `.github/workflows/deploy-ec2.yml` (one secret per stack).

---

## Files in this directory

| Path | Role |
|------|------|
| `mail/docker-compose.yml` | Mailpit service + shared network. |
| `database/docker-compose.yml` | MySQL 8 on `helpdesk-shared`; init SQL creates `secure_helpdesk` and grants `admin` access to both app databases. |
| `vuln/docker-compose.yml` | Vulnerable variant: backend + web images and env. |
| `secure/docker-compose.yml` | Secure variant: backend + web images and env. |
| `nginx/chisomejim.site.conf` | Example Nginx server blocks for host → localhost port mapping. On EC2, TLS is configured with **Certbot**; certificate paths in the live config follow `/etc/letsencrypt/`. |
