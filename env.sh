for FOLDER in api-auth api-scrapper api-search tmdb-proxy client pg-admin postgres
do
    echo "source $FOLDER"
    export $(grep -v '^#' ./$FOLDER/.env | xargs)
done