FROM golang:1.24.4-alpine3.22

WORKDIR /app
COPY . .
RUN go build -o app ./cmd/server

CMD [ "./app" ]
