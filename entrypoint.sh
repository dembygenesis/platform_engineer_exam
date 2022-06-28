wait-for "${DB_HOST}:${DB_PORT}" -- "$@"

echo "#######################"
printenv
echo "#######################"

populate_admin
api_server