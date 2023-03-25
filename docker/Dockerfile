FROM postgres:latest

# Set environment variables for database
ENV POSTGRES_USER=postgres
ENV POSTGRES_PASSWORD=development
ENV POSTGRES_DB=kedaiprogrammer

# Copy initialization script
COPY docker/init.sql /docker-entrypoint-initdb.d/