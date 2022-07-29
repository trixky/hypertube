FROM node:16-alpine3.15 as builder

RUN yarn set version berry

WORKDIR /app
COPY . .

RUN yarn install
RUN yarn build

FROM node:16-alpine3.15

WORKDIR /app
COPY --from=builder --chown=node:node /app .

RUN mkdir -p /app/.cache
RUN chown -R node:node /app/.cache
USER node
ENTRYPOINT dumb-init node ./dist/index.js
