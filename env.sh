for FOLDER in api-auth api-user api-scrapper api-media tmdb-proxy api-streaming api-position client postgres redis
do
    echo "source $FOLDER"
    export $(grep -v '^#' ./$FOLDER/.env | xargs)
done