FROM alpine:latest as postgresql
RUN apk add --no-cache postgresql-client
WORKDIR /app
COPY ./docker-entrypoint-initdb.d /docker-entrypoint-initdb.d/
ENV POSTGRES_USER=postgres
ENV POSTGRES_PASSWORD=development
ENV POSTGRES_DB=kedaiprogrammer
RUN mkdir /run/postgresql && chown -R postgres:postgres /run/postgresql
USER postgres
RUN initdb
RUN pg_ctl -D /var/lib/postgresql/data -l logfile start &&\
    psql --command "CREATE USER $POSTGRES_USER WITH SUPERUSER PASSWORD '$POSTGRES_PASSWORD';" &&\
    createdb -O $POSTGRES_USER $POSTGRES_DB &&\
    pg_ctl -D /var/lib/postgresql/data -l logfile stop

FROM golang:1.20 as build
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o main .

FROM alpine:latest
RUN apk add --no-cache ca-certificates
COPY --from=postgresql /usr/bin/pg_ctl /usr/bin/
COPY --from=postgresql /usr/bin/postgres /usr/bin/
COPY --from=postgresql /usr/bin/initdb /usr/bin/
COPY --from=postgresql /usr/lib/postgresql/ /usr/lib/postgresql/
COPY --from=postgresql /usr/share/postgresql/ /usr/share/postgresql/
COPY --from=postgresql /etc/postgresql/ /etc/postgresql/
COPY --from=postgresql --chown=postgres:postgres /run/postgresql/ /run/postgresql/
COPY --from=build /app/main /app/
ENV DB_HOST=localhost
ENV DB_PORT=5432
ENV DB_NAME=kedaiprogrammer
ENV DB_USER=postgres
ENV DB_PASSWORD=development
EXPOSE 8080
CMD ["/app/main", "--host", "0.0.0.0"]
