wait-for "${DB_HOST}:${DB_PORT}" -- "$@"

echo "#######################"
printenv
echo "#######################"

/usr/bin/populate_admin
/usr/bin/api_server