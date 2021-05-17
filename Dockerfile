FROM golang:1.16

WORKDIR /go
COPY main.go /go/main.go

RUN go build /go/main.go

ENTRYPOINT ["/go/main"]
