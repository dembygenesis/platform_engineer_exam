rm -rf ./src/persistence/mysql/models_schema
mkdir -p ./src/persistence/mysql/models_schema -- "$1"

sqlboiler mysql -p models_schema --no-tests