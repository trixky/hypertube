FROM node:16-alpine3.15 as builder

RUN yarn set version berry

WORKDIR /app
COPY . .

RUN yarn install
RUN yarn build

FROM builder

WORKDIR /app
COPY --from=builder /app .

ENTRYPOINT PORT=${CLIENT_PORT} node ./build/index.js
