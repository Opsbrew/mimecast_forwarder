FROM golang:1.16

WORKDIR /go/src/github.com/opsbrew/mimecast_forwarder

COPY . .

RUN go mod init github.com/opsbrew/mimecast_forwarder | true
RUN go mod vendor
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor
RUN ls

FROM alpine

WORKDIR /app

COPY --from=0 /go/src/github.com/opsbrew/mimecast_forwarder .

RUN addgroup -S myawesomegroup
RUN adduser -S myawesomeuser -G myawesomegroup
USER myawesomeuser

CMD [ "./mimecast_forwarder","start" ]
# CMD [ "sleep","3000" ]