FROM golang:1.8-alpine
RUN apk update && apk upgrade && apk add --no-cache git

RUN mkdir -p /go/src/github.com/slonoed/chantra/
ADD . /go/src/github.com/slonoed/chantra
WORKDIR /go/src/github.com/slonoed/chantra

RUN go build
RUN go install

EXPOSE 3000

CMD ["chantra"]
