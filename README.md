# Gonertia

Goravel Lite, PostgreSQL, Inertia, Vue 3, and TypeScript behind Traefik.
Includes two test pages; no UI libraries, authentication, or SSR yet.

## Setup

Requires Docker Compose, VS Code Dev Containers, and a running Traefik instance
on the `traefik_proxy` network with an HTTP entrypoint named `web`.
The host `.ssh`, `.gitconfig`, and OpenCode configuration/state directories
referenced in `.devcontainer/devcontainer.json` must exist.

1. Copy `.env.example` to `.env`; optionally adjust domains or database credentials.
2. Choose **Dev Containers: Reopen in Container**.
3. In the container terminal, from `/workspace`, run once:

```sh
chmod 600 .env
go mod download
npm ci
./artisan key:generate
./artisan migrate
```

Keep `.env` private. Generate the app key only once per checkout.

Without VS Code:

```sh
docker compose -f docker-compose.dev.yml up -d
docker compose -f docker-compose.dev.yml exec app bash
```

## Development

Run manually in two container terminals:

```sh
air          # Terminal 1: Goravel
npm run dev  # Terminal 2: Vite
```

Open **http://goravel.localhost**. `/` and `/about` demonstrate Inertia navigation.
Vite assets/HMR use **http://vite.goravel.localhost**. Press **Ctrl+C** in each
terminal to stop its server; neither starts automatically.

- Traefik handles ports `3000` and `5173`; VS Code auto-forwarding is disabled.
- After changing Compose settings or domains, use **Dev Containers: Rebuild Container**.
  Stop any old forwards in VS Code's **Ports** panel.
- For multiple projects, use unique `APP_DOMAIN` / `VITE_DOMAIN` values.
- To use built assets, stop Vite, run `npm run build`, and restart Air.
  If Vite was killed abruptly, remove stale `public/hot` first.

Stop containers without deleting database data:

```sh
docker compose -f docker-compose.dev.yml down
```

## Checks

```sh
npm run typecheck
npm run test:config
npm run build
go test ./...
go vet ./...
```

With Goravel running, check database readiness from the host:

```sh
curl --fail http://goravel.localhost/health/ready
```

Optional PostgreSQL integration test, using a separate disposable database:

```sh
# Create once, inside the container.
PGPASSWORD="$DB_PASSWORD" createdb -h "$DB_HOST" -U "$DB_USERNAME" gonertia_test_smoke
GONERTIA_TEST_DATABASE=gonertia_test_smoke go test -count=1 ./...
```

## Optional agent tooling

`.agents/` is ignored by Git. Install the official Goravel skill per checkout:

```sh
mkdir -p .agents/skills/goravel-development
curl --fail --location \
  https://raw.githubusercontent.com/goravel/goravel-lite/20a55f95b6ebca8411a81a687054fa944ed55939/.agents/skills/goravel-development/SKILL.md \
  --output .agents/skills/goravel-development/SKILL.md
```

Review the skill, then start a new OpenCode session. It is discovered
automatically; mention `@goravel-development` to load it explicitly.
Add local conventions in the same directory's `CUSTOM.md` if needed.

The Nuxt UI MCP is configured in `opencode.json`; check with `opencode mcp list`.
