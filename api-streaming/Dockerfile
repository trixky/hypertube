FROM node:16-alpine3.15 as builder

# install ffmpeg for transcode
RUN apk add --no-cache ffmpeg

RUN yarn set version stable

WORKDIR /app
COPY ./package.json .
COPY ./.yarnrc.yml .
COPY ./yarn.lock .

RUN yarn install

FROM builder

WORKDIR /app

COPY . .

ENTRYPOINT npx ts-node-dev --respawn main.ts
