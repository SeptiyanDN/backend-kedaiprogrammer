FROM postgres:15.1-alpine

# Set environment variables for database
ENV POSTGRES_USER=postgres
ENV POSTGRES_PASSWORD=development
ENV POSTGRES_DB=kedaiprogrammer

# Copy initialization script


LABEL author="Your Name"
LABEL description="Postgres Image for demo"
LABEL version="1.0"

COPY *.sql /docker-entrypoint-initdb.d/