# https://nodejs.org/en/docs/guides/nodejs-docker-webapp/
FROM node:16-alpine3.15

# install ffmpeg for convert mkv to mp4
RUN apk add --no-cache ffmpeg

# copy api files
WORKDIR /app
COPY ./package.json .

# install CompileDaemon (hot reload)
RUN yarn install
ENTRYPOINT npx ts-node-dev --respawn main.ts
