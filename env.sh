for FOLDER in postgres redis api-auth api-user api-scrapper api-media tmdb-proxy api-streaming api-position client

do
    echo "source $FOLDER"
    export $(grep -v '^#' ./$FOLDER/.env | xargs)
done