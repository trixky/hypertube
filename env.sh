#!/bin/bash

# server / postgres
export POSTGRES_USER=hyper
export POSTGRES_PASSWORD=1234
export POSTGRES_DB=hypertube
export POSTGRES_HOST=postgres
export POSTGRES_PORT=5432

# pgadmin
export PGADMIN_DEFAULT_EMAIL=admin@hypertube.com
export PGADMIN_DEFAULT_PASSWORD=1234

# supertest / postgres
export PGUSER=$POSTGRES_USER
export PGPASSWORD=$POSTGRES_PASSWORD
export PGDATABASE=$POSTGRES_DB

# supertest
export PGHOST=$POSTGRES_HOST
export PGPORT=$POSTGRES_PORT