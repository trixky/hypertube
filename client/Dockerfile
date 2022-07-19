FROM node:16-alpine3.15

WORKDIR /app
COPY ./package.json .
# COPY ./yarn.lock .

RUN yarn install
COPY . .

ENTRYPOINT yarn run dev --host --port ${CLIENT_PORT}
