echo \#\#\# start postgres/redis

# Start postgres
/usr/local/bin/docker-entrypoint.sh postgres &> /dev/null &
# Start redis
redis-server &>/dev/null &

# Wait for postgres and redis to be started
sleep 5

# Change the network context to localhost
export POSTGRES_HOST=localhost
export REDIS_HOST=localhost

echo ============================================= START TESTS

# for each APIs
for API in api-auth api-user api-media api-scrapper
do
    cd /hypertube/$API
    echo --------------------- TEST '['$API']'
    echo \#\#\# download dependencies
    # Download missing depedencies
    go get -d ./... &>/dev/null
    # Start tests
    # Ignore some generated folder
    go test $(go list ./... | grep -v /proto | grep -v /sqlc)
done
