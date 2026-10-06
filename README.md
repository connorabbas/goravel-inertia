# Gonertia

Goravel Lite 1.18 with Gin, PostgreSQL, and an existing Traefik proxy.
Step 2 implements the backend/database foundation only. Inertia, Vue,
TypeScript, Nuxt UI integration, and authentication are deferred.

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
- `http://vite.goravel.localhost` to Vite on port 5173 (reserved for the next step)

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

The landing response is `Gonertia: Goravel is running.` Readiness executes
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

HTTP tests cover the landing page and readiness success/failure without needing
a live database. The PostgreSQL integration test is opt-in. Create a separate,
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
- Developers start and stop `air` manually; agents should not start it unasked.
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
