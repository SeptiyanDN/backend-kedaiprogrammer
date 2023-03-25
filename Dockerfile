FROM golang:1.20

# Set the current working directory inside the container
# WORKDIR /app

# Copy the Go modules and install them
COPY go.mod go.sum ./
RUN go mod tidy

# Copy the rest of the application code
COPY . .

# Build the application
RUN go build -o main .

# Use a minimal image of Alpine for the final image
FROM alpine:latest

# Add PostgreSQL as a dependency
RUN apk add --no-cache postgresql-client

# Set the current working directory inside the container
# WORKDIR /app

# Copy the built binary from the previous stage
# COPY --from=0 /app/main .

ENV DB_HOST=localhost
ENV DB_PORT=5432
ENV DB_NAME=kedaiprogrammer
ENV DB_USER=postgres
ENV DB_PASSWORD=development

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["./main", "--host", "0.0.0.0"]
