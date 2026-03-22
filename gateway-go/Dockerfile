FROM golang:1.24-alpine

WORKDIR /workspace
COPY . .

WORKDIR /workspace/gateway-go
RUN go build -o /bin/gateway ./cmd/gateway

EXPOSE 8080
CMD ["/bin/gateway"]
