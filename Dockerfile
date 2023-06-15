FROM golang:1.19-alpine as builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o linux-web-top main.go

FROM alpine
COPY --from=builder /app/linux-web-top /usr/local/bin/linux-web-top
EXPOSE 3000

ENTRYPOINT ["linux-web-top" , "-l=:3000"]