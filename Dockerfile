FROM golang:alpine AS builder
RUN apk update && apk add --no-cache git
WORKDIR $GOPATH/src/Comonut/vectorugo/
COPY . .
RUN go get -d -v
RUN CGO_ENABLED=0 go build -o /go/bin/main
FROM scratch
COPY --from=builder /go/bin/main /go/bin/main
ENTRYPOINT ["/go/bin/main"]
