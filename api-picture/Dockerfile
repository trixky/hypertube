FROM node:16-alpine3.15 as builder

RUN yarn set version berry

WORKDIR /app
COPY . .

RUN yarn install
RUN yarn build

FROM node:16-alpine3.15

RUN apk add dumb-init

WORKDIR /app
COPY --from=builder --chown=node:node /app .

RUN mkdir -p /app/storage
RUN chown -R node:node /app/storage
USER node
ENTRYPOINT dumb-init node ./dist/index.js
