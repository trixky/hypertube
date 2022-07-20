# https://nodejs.org/en/docs/guides/nodejs-docker-webapp/

FROM node:16-alpine3.15

WORKDIR /usr/src/app

# copy configuration and project files
COPY . .

# install project
RUN npm install

ENTRYPOINT npx jest --runInBand