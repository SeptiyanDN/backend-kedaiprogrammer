# Use an official Golang runtime as a parent image
FROM golang:1.20

# Set the working directory to /app
WORKDIR /app

# Copy the current directory contents into the container at /app
COPY . /app

# Build the application
RUN go build -o main .

# Use an official Alpine Linux runtime as a parent image
FROM alpine:latest

# Add PostgreSQL as a dependency
RUN apk add --no-cache postgresql-client

# Set the working directory to /app
WORKDIR /app

# Copy the binary from the previous stage
COPY --from=0 /app/main .

# Set environment variables for database connection
ENV DB_HOST=localhost
ENV DB_PORT=5432
ENV DB_NAME=kedaiprogrammer
ENV DB_USER=postgres
ENV DB_PASSWORD=development

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["./main", "--host", "0.0.0.0"]
