FROM golang:1.24-alpine

WORKDIR /app
COPY . .

RUN go build -o /bin/gateway ./cmd/gateway

EXPOSE 8080
CMD ["/bin/gateway"]
