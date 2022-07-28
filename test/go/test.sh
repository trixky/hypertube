echo \#\#\# start postgres/redis

# start postgres
/usr/local/bin/docker-entrypoint.sh postgres &> /dev/null &
# start redis
redis-server &>/dev/null &

# wait for postgres and redis to be started
sleep 5

# change the network context to localhost
export POSTGRES_HOST=localhost
export REDIS_HOST=localhost

echo ============================================= START TESTS

# for each API
for API in api-auth api-user
do
    cd /hypertube/$API
    echo --------------------- TEST '['$API']'
    echo \#\#\# download dependencies
    # download missing depedencies
    go get -d ./... &>/dev/null
    # start test
    go test ./...
done
