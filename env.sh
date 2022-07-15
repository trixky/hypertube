for FOLDER in api-auth api-user api-scrapper api-media tmdb-proxy api-streaming client postgres
do
    echo "source $FOLDER"
    export $(grep -v '^#' ./$FOLDER/.env | xargs)
done