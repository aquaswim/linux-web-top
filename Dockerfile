FROM golang:1.19-alpine as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o linux-web-top main.go

FROM alpine
WORKDIR /app
COPY --from=builder /app/linux-web-top linux-web-top
COPY web web
EXPOSE 3000

ENTRYPOINT ["./linux-web-top" , "-l=:3000", "-p=/hostproc"]
