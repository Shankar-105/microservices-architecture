FROM golang:1.24-alpine

WORKDIR /workspace
COPY . .

WORKDIR /workspace/user-service-go
RUN go build -o /bin/user-service ./cmd/user

EXPOSE 50051
CMD ["/bin/user-service"]
