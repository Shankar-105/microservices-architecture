FROM golang:1.24-alpine

WORKDIR /app
COPY . .

RUN go build -o /bin/user-service ./cmd/user

EXPOSE 50051
CMD ["/bin/user-service"]
