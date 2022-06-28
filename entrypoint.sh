wait-for "${DB_HOST}:${DB_PORT}" -- "$@"

populate_admin
api_server