FROM golang:1.21-alpine

WORKDIR /app

COPY . .

RUN go build -o stress-tester .

ENTRYPOINT ["./stress-tester"]
