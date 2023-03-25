# Use a minimal image of Alpine for the final image
FROM alpine:latest

# Add PostgreSQL as a dependency
RUN apk add --no-cache postgresql-client

# Set the current working directory inside the container
WORKDIR /app

FROM golang:1.20

# Set the current working directory inside the container
WORKDIR /app
# Copy the rest of the application code
COPY . .

RUN go mod tidy

RUN go build -o main .

# Set environment variables for database connection
ENV DB_HOST=localhost
ENV DB_PORT=5432
ENV DB_NAME=kedaiprogrammer
ENV DB_USER=postgres
ENV DB_PASSWORD=development

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["/app/main", "--host", "0.0.0.0"]
