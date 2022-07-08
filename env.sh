for FOLDER in api-auth api-user  client pg-admin postgres
do
    echo "source $FOLDER"
    export $(grep -v '^#' ./$FOLDER/.env | xargs)
done