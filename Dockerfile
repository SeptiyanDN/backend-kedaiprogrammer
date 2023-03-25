FROM golang:1.20-alpine AS build
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o main .

# Final stage
FROM alpine:latest
RUN apk add --no-cache postgresql-client
WORKDIR /app
COPY --from=build /app/main .
ENV DB_HOST=localhost
ENV DB_PORT=5432
ENV DB_NAME=kedaiprogrammer
ENV DB_USER=postgres
ENV DB_PASSWORD=development
EXPOSE 8080
CMD ["/app/main", "--host", "0.0.0.0"]

