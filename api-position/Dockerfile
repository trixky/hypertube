# https://nodejs.org/en/docs/guides/nodejs-docker-webapp/
FROM node:16-alpine3.15

WORKDIR /app
COPY ./package.json .
# COPY ./yarn.lock .

RUN yarn install
COPY . .

ENTRYPOINT npx ts-node-dev --respawn main.ts
