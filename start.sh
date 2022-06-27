docker container rm platform_engineer_db -f
docker container rm platform_engineer_api -f
docker-compose -f  docker-compose.yml down --remove-orphans --volumes
docker-compose up --force-recreate --build