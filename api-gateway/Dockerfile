FROM golang:1.25.1-alpine

WORKDIR /api-gateway-app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main .

CMD ["./main"]
