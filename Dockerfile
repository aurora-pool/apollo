FROM golang:1.10

RUN go get -u github.com/kardianos/govendor

RUN mkdir -p /go/src/github.com/aurora-pool/apollo
WORKDIR /go/src/github.com/aurora-pool/apollo

ADD . /go/src/github.com/aurora-pool/apollo

RUN govendor sync && go install github.com/aurora-pool/apollo

EXPOSE 8242
