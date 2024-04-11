FROM mariadb:10.4

CMD echo $DB_DATABASE
COPY ./database/*.sql /docker-entrypoint-initdb.d/

RUN apt update && apt install nano
RUN chmod 755 /docker-entrypoint-initdb.d/*.sql