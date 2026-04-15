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

1. **Backend configuration** — The workflow loads environment variables from **AWS Secrets Manager** (different secret per stack). Those values are used when building the backend image so settings like database credentials are baked in for that variant.

2. **Build & push images** — Backend and frontend **Docker** images are built in CI, tagged `vuln` or `secure`, plus the git commit SHA, and pushed to **Amazon ECR Public** so the server can pull them by name.

3. **Frontend API URL** — The frontend is built with a fixed **HTTPS** API URL for that stack.

4. **Deploy to EC2 over SSH** — The workflow connects to the instance using SSH private key and server IP in  **repository secret**. It copies the **Compose** files from this repo and runs `docker compose pull` and `up` so the running containers match the new images.

---

## Domains and URLs (as configured in this repo)

The workflow and Compose files assume a fixed production domain (`chisomejim.site`) with subdomains for:

| Resource | URL |
|----------|-----|
| Vuln (UI) | https://vulnweb.chisomejim.site |
| Vuln (API) | https://api.vuln.chisomejim.site |
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
| `~/helpdesk/vuln/` or `~/helpdesk/secure/` | The **app stack** for that variant (backend + web containers). |

The **mail** stack is updated on every deploy. The **app** stack for the branch you pushed is updated the same way.

---

## Docker networking

- The **mail** Compose file defines a shared Docker network called **`helpdesk-shared`**.
- Each **app** stack joins that network as **external** so backends can reach Mailpit by hostname (`mailpit`) on the internal SMTP port.
- App HTTP ports are bound to **`127.0.0.1` only** on the host. They are not meant to be reached directly from the internet; traffic is expected to come through a **reverse proxy** on the same machine (see below).

---

## Reverse proxy (Nginx)

`nginx/chisomejim.site.conf` is an example **Nginx** configuration: it maps public hostnames to those localhost ports (separate hosts for each UI, each API, and Mailpit’s web UI). The Docker only listens on the loopback interface.

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
| `vuln/docker-compose.yml` | Vulnerable variant: backend + web images and env. |
| `secure/docker-compose.yml` | Secure variant: backend + web images and env. |
| `nginx/chisomejim.site.conf` | Example Nginx server blocks for host → localhost port mapping. |
