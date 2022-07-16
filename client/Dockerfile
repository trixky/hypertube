FROM node:16-alpine3.15

# https://sveltesociety.dev/recipes/publishing-and-deploying/dockerize-a-svelte-app/

# copy api files
WORKDIR /app
COPY ./package.json .

# install CompileDaemon (hot reload)
RUN yarn install
ENTRYPOINT yarn run dev --host --port ${CLIENT_PORT}
