FROM golang:1.24-alpine

WORKDIR /workspace
COPY . .

WORKDIR /workspace/order-service-go
RUN go build -o /bin/order-service ./cmd/order

EXPOSE 50052
CMD ["/bin/order-service"]
