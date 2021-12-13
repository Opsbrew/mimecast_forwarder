FROM golang:1.17

WORKDIR /go/src/github.com/opsbrew/mimecast_forwarder

COPY . .

RUN go mod init github.com/opsbrew/mimecast_forwarder | true
RUN go mod vendor
RUN go get github.com/cooldrip/cstrftime
RUN CGO_ENABLED=0 GOOS=linux go build main.go
RUN ls
FROM alpine

WORKDIR /app

COPY --from=0 /go/src/github.com/opsbrew/mimecast_forwarder .

RUN addgroup -S myawesomegroup
RUN adduser -S myawesomeuser -G myawesomegroup
USER myawesomeuser

EXPOSE 8080
ENTRYPOINT [ "./mimecast_forwarder" ]