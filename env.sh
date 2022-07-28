echo "source root"
if test -f "./.env"; then
    export $(grep -v '^#' ./.env | xargs)
else
    echo "Missing root .env file"
fi

for FOLDER in postgres redis api-auth api-user api-scrapper api-media tmdb-proxy api-streaming api-position client streaming-proxy
do
    echo "source $FOLDER"
    if test -f "./$FOLDER/.env"; then
        export $(echo "$(echo $(cat ./$FOLDER/.env) | grep -v '^#' | xargs)") 2>/dev/null
    else
        echo "No .env file in $FOLDER"
    fi
done
