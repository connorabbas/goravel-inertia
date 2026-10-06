# Goravel Inertia starter

Development environment for a Goravel application with Inertia, Vue,
TypeScript, Nuxt UI, PostgreSQL, and an existing Traefik proxy.

## Prerequisites

- Docker with Compose
- A running Traefik instance attached to an external Docker network named
  `traefik_proxy`, with an HTTP entrypoint named `web`
- VS Code with Dev Containers, or another compatible Dev Container client

If the proxy network does not exist yet, create it once:

```sh
docker network create traefik_proxy
```

Optionally copy `.env.example` to `.env` to change domains or database
credentials. The defaults work without an `.env` file.

## Open the development container

Choose **Dev Containers: Reopen in Container**. This builds the development
image and starts PostgreSQL. The application container intentionally runs an
idle command until the Goravel project is scaffolded.

The environment provides Go 1.26, Node.js 22, PostgreSQL's client, Air, and
OpenCode V2. Run `opencode` in the container terminal. The Dev Container shares
your host OpenCode configuration, sessions, and state; those host directories
and your `.gitconfig` and `.ssh` must exist before reopening the container.
Dependency and database data are held in named Docker volumes.

Once the application servers are running, Traefik routes:

- `http://goravel.localhost` to Goravel on port 3000
- `http://vite.goravel.localhost` to Vite on port 5173

Goravel should listen on `0.0.0.0:3000`, and Vite should listen on
`0.0.0.0:5173`. The container already supplies the matching application,
database, and Vite environment variables.

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
