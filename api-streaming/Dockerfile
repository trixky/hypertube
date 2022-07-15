# https://nodejs.org/en/docs/guides/nodejs-docker-webapp/

FROM node:16-alpine3.15

WORKDIR /usr/src/app

# install ffmpeg for convert mkv to mp4
RUN apk add --no-cache ffmpeg

# copy configuration files
COPY *.json ./

# install project
RUN npm install

# copy project
COPY . .

ENTRYPOINT npx ts-node-dev --respawn main.ts