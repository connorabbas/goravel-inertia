# Gonertia

Goravel Lite 1.18 with Gin, PostgreSQL, Inertia 3, Vue 3, TypeScript, and Vite
behind an existing Traefik proxy. Two plain pages demonstrate server-driven
routing; UI libraries, authentication, and SSR are not installed.

Scaffold: `goravel/goravel-lite` revision
`20a55f95b6ebca8411a81a687054fa944ed55939`. Framework and drivers are pinned
in `go.mod`; dependency checksums are recorded in `go.sum`.

## Prerequisites

- Docker with Compose
- A running Traefik instance attached to an external Docker network named
  `traefik_proxy`, with an HTTP entrypoint named `web`
- VS Code with Dev Containers, or another compatible Dev Container client

If the proxy network does not exist yet, create it once:

```sh
docker network create traefik_proxy
```

Copy `.env.example` to `.env` before starting Compose if you want to change
domains or database credentials. Each checkout needs its own generated app key.

## Open the development container

Choose **Dev Containers: Reopen in Container**. This builds the development
image and starts PostgreSQL. The application container intentionally runs an
idle command; start Goravel from its terminal using the commands below.

The environment provides Go 1.26, Node.js 22, PostgreSQL's client, Air, and
OpenCode V2. Run `opencode` in the container terminal. The Dev Container shares
your host OpenCode configuration, sessions, and state; those host directories
and your `.gitconfig` and `.ssh` must exist before reopening the container.
Dependency and database data are held in named Docker volumes.

Once the application servers are running, Traefik routes:

- `http://goravel.localhost` to Goravel on port 3000
- `http://vite.goravel.localhost` to Vite assets and HMR on port 5173

Goravel should listen on `0.0.0.0:3000`, and Vite should listen on
`0.0.0.0:5173`. The container already supplies the matching application,
database, and Vite environment variables.

## Bootstrap and start Goravel

Run inside the container, from `/workspace`:

```sh
cp -n .env.example .env
chmod 600 .env
go mod download
./artisan key:generate
./artisan migrate
air
```

Generate the key only once per checkout; do not regenerate it at every startup.
For normal development, open a dev-container terminal in `/workspace` and run
`air` manually. Keep that terminal open; press **Ctrl+C** to stop Air and its
Goravel server. Run `air` again to restart. Neither Compose nor the dev container
starts Air automatically.

The `.env` file is ignored by Git. Compose supplies the HTTP and database
settings, which take precedence over `.env`; custom Compose values require
recreating the app container. Restart Air after changing `.env`.

`DB_SSLMODE=disable` is intended for local PostgreSQL only; configure TLS for
remote/production databases. PostgreSQL data persists in its named volume;
changing credentials in `.env` does not update an already-initialized database.

There are no application migrations yet. `migrate` initializes the migration
repository; `migrate:status` reports no migrations until you add one with
`./artisan make:migration`. Migrations do not run automatically on startup.

## Verify

From the host, with Air running:

```sh
curl --fail http://goravel.localhost/
curl --fail http://goravel.localhost/health/ready
```

The landing response is now an HTML shell mounting the Vue `Home` page through
Inertia. Start Vite or build the frontend as described below before opening it
in your browser. Readiness executes
`SELECT 1` through Goravel with a two-second query deadline and returns
`{"status":"ok"}`. Database query failures return HTTP 503 with a generic
response, without exposing credentials. For a direct container check, use
`http://127.0.0.1:3000/health/ready`.

Inside the container:

```sh
go build -o tmp/gonertia .
go test ./...
go vet ./...
./artisan migrate:status
```

HTTP tests cover both pages' HTML and Inertia JSON responses, backend props,
stale asset versions, and readiness success/failure without needing a live
database. After a frontend build, they also verify static asset routing.
The PostgreSQL integration test is opt-in. Create a separate,
disposable test database once:

```sh
PGPASSWORD="$DB_PASSWORD" createdb -h "$DB_HOST" -U "$DB_USERNAME" gonertia_test_smoke
```

Then run this repeatable check:

```sh
GONERTIA_TEST_DATABASE=gonertia_test_smoke go test -count=1 -v ./...
```

The test selects its database before providers boot, verifies the database
name and readiness endpoint, migrates a test-only table, inserts/reads a row,
checks migration status, and rolls back. It rejects the configured development
database and names without the `gonertia_test_` prefix. Do not point it at a
database containing valuable data. The database and empty migration repository
remain available for subsequent runs.

Air watches Go source changes, including migrations, and ignores test files,
runtime storage, build output, Git, and future frontend dependencies.

## Vue / Inertia development

Install the frontend dependencies once per checkout, inside the dev container:

```sh
npm ci
```

Use two terminals in `/workspace`, starting both servers yourself:

```sh
# Terminal 1: Go backend with reload
air
```

```sh
# Terminal 2: Vue / TypeScript assets with HMR
npm run dev
```

Open **http://goravel.localhost**, not the Vite domain. `/` and `/about` are
Goravel routes; Inertia's `<Link>` visits them without a full document reload.
There is no Vue Router or separate frontend API. Props are passed by the Go
controllers and typed using Vue's `<script setup lang="ts">`.

Press **Ctrl+C** in each terminal to stop that server. Neither server starts
automatically. Vite writes its browser-accessible origin into ignored
`public/hot` while running and removes it on normal shutdown. If Vite is killed
abruptly, remove the stale file with `rm -f public/hot` before using built assets.

Vite always listens on `0.0.0.0:5173`; assets and WebSocket HMR are proxied through
`http://vite.goravel.localhost` on browser port 80. Allowed hosts and CORS are
restricted to the configured Vite hostname and `APP_URL`. Compose supplies
`VITE_ORIGIN` and optional `VITE_USE_POLLING`;
customize `APP_DOMAIN` / `VITE_DOMAIN` in `.env` and recreate the app container
when changing them. The old `VITE_DEV_URL` variable is not used: only
`public/hot` enables development assets, so built assets work without Vite.

VS Code automatic forwarding is disabled for ports `3000` and `5173` using
`portsAttributes` in `.devcontainer/devcontainer.json`; other ports can still
be forwarded. Stop any existing forwards in VS Code's **Ports** panel. Use
**Dev Containers: Rebuild Container** once to apply the updated container
configuration and Compose environment; a cache-clearing rebuild is unnecessary.
Parallel projects can reuse the internal ports with unique `APP_DOMAIN` and
`VITE_DOMAIN` values, rather than allocating different host ports.

### Frontend checks and built assets

```sh
npm run typecheck
npm run test:config
npm run build
go test ./...
```

The build writes a manifest and hashed assets into ignored `public/build`.
`test:config` checks Traefik/HMR settings and the hot-file lifecycle in a
temporary directory without starting Vite or opening a port.
With Vite stopped and no `public/hot`, Goravel loads `/build/` assets from this
manifest. Restart Air after rebuilding: the adapter caches the manifest and
asset version during application boot. `vite preview` is not needed; Goravel
serves the app in both modes.

Manual browser checks:

- Load `/` and `/about` directly; both show messages supplied by Go.
- Follow the Inertia links; navigation requests carry `X-Inertia: true` and
  return page JSON. Browser back/forward should work.
- Edit a `.vue` template while Vite is running and confirm HMR updates it.
- Stop Vite, run `npm run build`, restart Air, and confirm both pages work
  without the Vite server.

## Agent tooling

### Install the Goravel development skill

`.agents/` is ignored by Git, so each new checkout needs to install the skill
locally. From the project root inside the dev container, run:

```sh
mkdir -p .agents/skills/goravel-development
curl --fail --location \
  https://raw.githubusercontent.com/goravel/goravel-lite/20a55f95b6ebca8411a81a687054fa944ed55939/.agents/skills/goravel-development/SKILL.md \
  --output .agents/skills/goravel-development/SKILL.md
```

This downloads the official skill from the same revision as this project's
scaffold. Review the downloaded instructions, then start a new OpenCode session
from the project root. OpenCode discovers `.agents/skills/` automatically;
mention `@goravel-development` to explicitly load the skill. No extra
`opencode.json` configuration is needed.

Optionally create `.agents/skills/goravel-development/CUSTOM.md` with local
project conventions; the official skill instructs agents to read it:

```markdown
# Gonertia conventions

- Run commands from `/workspace` inside the dev container.
- Use the standard Goravel layout and the `goravel` module name.
- Use Gin on `0.0.0.0:3000` and PostgreSQL at `postgres:5432`.
- Preserve the existing Docker/Traefik setup.
- Developers start and stop `air` and `npm run dev` manually; agents should not start them unasked.
- Never commit `.env` or generated keys; Compose variables override `.env`.
- Database write tests must use a separate `gonertia_test_*` database.
- Verify changes with `go test ./...` and `go vet ./...`.
```

Both files remain local and untracked. When upgrading Goravel, review the
matching upstream skill and update `SKILL.md`, preserving your `CUSTOM.md`.

### Nuxt UI MCP

`opencode.json` configures the existing Nuxt UI MCP using native V2 syntax;
verify with `opencode mcp list`. No additional MCP servers or Goravel runtime
AI packages are required.

## Terminal-only workflow

Without a Dev Container client, start the same environment and open a shell:

```sh
docker compose -f docker-compose.dev.yml up -d
docker compose -f docker-compose.dev.yml exec app bash
```

Stop it without deleting PostgreSQL data:

```sh
docker compose -f docker-compose.dev.yml down
```
