FROM golang:1.23 AS builder

ENV GOBIN /go/bin
ENV GO111MODULE=on

COPY ./go.mod ./go.sum /app/

WORKDIR /app

RUN go mod download

COPY . /app/

RUN CGO_ENABLED=0 go build -o main .

# Stage 2

FROM alpine:3.16 AS production

COPY --from=builder /app/main /app/

WORKDIR /app

EXPOSE 8080

CMD ["/app/main"]
