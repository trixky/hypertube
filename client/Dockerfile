FROM node:16-alpine3.15 as builder

RUN yarn set version stable

FROM builder

WORKDIR /app
COPY ./package.json .
COPY ./.yarnrc.yml .
COPY ./yarn.lock .

RUN yarn install
COPY . .

ENTRYPOINT yarn run dev --host --port ${CLIENT_PORT}
