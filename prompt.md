# Goravel with Vue Inertia Starter Kit

I want to create a new Goravel project and use their associated inertia package/adapter for a Vue frontend with TypeScript and Nuxt UI. The end goal is to have a fucntional starter kit similar to https://github.com/connorabbas/laravel-nuxtui-starter-kit but in Go, with a smaller scope for packages that dont have an equivilent (e.g. the spate data / TS transformer and laravel foritfy), we can have it be a simple authentication starter kit using the tooling provided by Goravel.

Reference docs: https://github.com/goravel/inertia 
https://www.goravel.dev/
https://github.com/goravel/goravel

## Step 1 - COMPLETED
The first step is to create the development environment. We are in the project root, I want to create a docker-compose development environment that includes the necessary services for the application and uses dev containers to provide a consistent development environment across different machines. Traefik will be used to proxy the traffic to the main container running the Goravel application, use the Laravel Nuxt UI starter kit as a refernce for a similar setup. Once the environment is set up, then we can start developing the Goravel application.

## Step 2
Now lets actually setup the goravel project and confirm everything works, only setting up goravel and connecting the database, ensuring the app works. Then we will add gonertia when everything is confirmed working end to end. Add any relevant MCP servers or agent skills for Go, goravel, inertia (frontend) that would be helpful before starting any real work and configure them with the opencode.json file or .agents/ dir, example: https://www.goravel.dev/ai/sdk.html#ai-agent-development-skill