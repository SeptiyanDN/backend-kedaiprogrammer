# Use an official Golang runtime as a parent image
FROM golang:1.20

# Set the working directory to /app
RUN apk add --no-cache postgresql-client

WORKDIR /app/main

# Copy the current directory contents into the container at /app
COPY go.mod .
COPY go.sum .
# Build the application
RUN go build -o main .

RUN go mod tidy

COPY . .
RUN go build -o ./out/main .


# Set environment variables for database connection
ENV DB_HOST=localhost
ENV DB_PORT=5432
ENV DB_NAME=kedaiprogrammer
ENV DB_USER=postgres
ENV DB_PASSWORD=development

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["./out/main", "--host", "0.0.0.0"]
