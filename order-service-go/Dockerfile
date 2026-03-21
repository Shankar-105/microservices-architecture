FROM golang:1.24-alpine

WORKDIR /app
COPY . .

RUN go build -o /bin/order-service ./cmd/order

EXPOSE 50052
CMD ["/bin/order-service"]
