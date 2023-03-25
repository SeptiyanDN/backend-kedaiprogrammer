FROM golang:1.20-alpine3.15

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

RUN go build -o main .

EXPOSE 8080

ENV DB_HOST=localhost
ENV DB_PORT=5432
ENV DB_NAME=kedaiprogrammer
ENV DB_USER=postgres
ENV DB_PASSWORD=development

CMD ["./main"]
