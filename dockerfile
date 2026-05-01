FROM golang:1.25-alpine AS builder 

WORKDIR /app

COPY go.mod go.sum ./ 
RUN go mod download 

COPY . . 
RUN go build -o url-shortener ./cmd/app

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/url-shortener .

EXPOSE 8080

ENTRYPOINT [ "./url-shortener" ]
CMD ["postgres"]