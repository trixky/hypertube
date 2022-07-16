FROM node:16-alpine3.15

# https://sveltesociety.dev/recipes/publishing-and-deploying/dockerize-a-svelte-app/

# copy api files
WORKDIR /app
COPY ./ /app

# install CompileDaemon (hot reload)
ENTRYPOINT npm install && npx svelte-kit dev --host --port ${CLIENT_PORT}
