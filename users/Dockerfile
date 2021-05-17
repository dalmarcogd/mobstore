FROM golang:1.16.4-buster as builder

WORKDIR /go/src/github.com/dalmarcogd/mobstore/users"/
COPY ./users" /go/src/github.com/dalmarcogd/mobstore/users"/
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GODEBUG=madvdontneed=1 go build -a -tags netgo -o application cmd/api/main.go && mv application /application

FROM alpine:3.12.0 as runner
WORKDIR /
COPY --from=builder /application .
EXPOSE 8080
ENTRYPOINT ["./application"]
