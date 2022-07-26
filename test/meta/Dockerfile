FROM node:16-alpine3.15 as builder

RUN yarn set version stable

WORKDIR /app
COPY ./package.json .
COPY ./.yarnrc.yml .
COPY ./yarn.lock .

RUN yarn install

FROM builder

WORKDIR /app

COPY . .

ENTRYPOINT yarn run jest --runInBand --forceExit
