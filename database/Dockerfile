FROM postgres

ENV POSTGRES_PASSWORD=${POSTGRES_PASSWORD}
ENV POSTGRES_USER=${POSTGRES_USER}
ENV POSTGRES_DB=${POSTGRES_DB}

RUN apt-get update && apt-get install -y postgresql-client
COPY ./migrations/db.sql /docker-entrypoint-initdb.d/
RUN chmod +x /docker-entrypoint-initdb.d/db.sql