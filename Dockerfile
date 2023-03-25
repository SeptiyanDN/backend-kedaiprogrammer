# Base image
FROM golang:1.20-alpine3.15

# Set the working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download and install go dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the app
RUN go build -o main .

# Expose port 8080
EXPOSE 8080

# Set environment variables
ENV DB_HOST=localhost
ENV DB_PORT=5432
ENV DB_NAME=kedaiprogrammer
ENV DB_USER=postgres
ENV DB_PASSWORD=development

# Start the app
CMD ["./main"]
