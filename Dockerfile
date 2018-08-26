FROM golang:1.10-alpine
WORKDIR /go/src/github.com/stephenhillier/leapi
ADD . /go/src/github.com/stephenhillier/leapi/
RUN go test && go install
