#!/bin/bash

# postgres
export POSTGRES_USER=hyper
export POSTGRES_PASSWORD=1234
export POSTGRES_DB=hypertube
export POSTGRES_HOST=postgres
export POSTGRES_PORT=5432
export PGUSER=$POSTGRES_USER
export PGPASSWORD=$POSTGRES_PASSWORD
export PGDATABASE=$POSTGRES_DB

# pgadmin
export PGADMIN_DEFAULT_EMAIL=admin@hypertube.com
export PGADMIN_DEFAULT_PASSWORD=1234

# # supertest
# export PGHOST=$POSTGRES_HOST
# export PGPORT=$

# api-auth
export API_AUTH_PORT=7010

# client
export CLIENT_PORT=4040