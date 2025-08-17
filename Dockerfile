FROM golang:1.25-alpine3.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -o /app/todo-api .

FROM scratch

COPY --from=builder /app/todo-api /app/todo-api

ENTRYPOINT ["/app/todo-api"]
