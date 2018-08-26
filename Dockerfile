FROM golang:1.10-alpine
RUN mkdir -p /go/src/github.com/stephenhillier/leapi
ADD . /go/src/github.com/stephenhillier/leapi/
RUN go install github.com/stephenhillier/leapi/
