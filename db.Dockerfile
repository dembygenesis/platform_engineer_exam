# FROM mysql:8.0.12
FROM mariadb:latest

COPY ./database/*.sql /docker-entrypoint-initdb.d/