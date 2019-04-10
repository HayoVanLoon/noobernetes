FROM golang:alpine AS builder

RUN apk add git
RUN apk add protobuf

WORKDIR /go/src/counter

COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

RUN ls /go/bin

FROM alpine
COPY --from=builder /go/bin/counter /usr/local/

CMD ["/usr/local/counter"]
