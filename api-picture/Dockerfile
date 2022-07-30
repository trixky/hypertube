FROM node:16-alpine3.15 as builder

RUN addgroup -S hypertube && adduser -S hypertube -G hypertube
RUN yarn set version berry

WORKDIR /app
COPY . .

RUN yarn install
RUN yarn build

FROM node:16-alpine3.15

WORKDIR /app
COPY --from=builder --chown=node:node /app .

RUN mkdir -p /app/storage
RUN chown -R hypertube:hypertube /app/storage
USER hypertube
ENTRYPOINT dumb-init node ./dist/index.js
