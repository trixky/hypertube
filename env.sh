echo "source root"
if test -f "./.env"; then
    export $(grep -v '^#' ./.env | xargs)
else
    echo "Missing root .env file"
fi

for FOLDER in postgres redis api-auth api-user api-picture api-scrapper api-media tmdb-proxy api-streaming api-position client streaming-proxy
do
    echo "source $FOLDER"
    if test -f "./$FOLDER/.env"; then
        if [ -n "$ZSH_VERSION" ]; then
            export $(grep -v '^#' ./$FOLDER/.env | xargs)
        else
            export $(echo "$(echo $(cat ./$FOLDER/.env) | grep -v '^#' | xargs)") 2>/dev/null
        fi
    else
        echo "No .env file in $FOLDER"
    fi
done
