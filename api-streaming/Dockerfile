# https://nodejs.org/en/docs/guides/nodejs-docker-webapp/
FROM node:16-alpine3.15 as builder

RUN yarn set version stable

FROM builder

# install ffmpeg for transcode
RUN apk add --no-cache ffmpeg

WORKDIR /app
COPY ./package.json .
COPY ./.yarnrc.yml .
COPY ./yarn.lock .

RUN yarn install
COPY . .

ENTRYPOINT npx ts-node-dev --respawn main.ts
