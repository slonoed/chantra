FROM golang:1.8-alpine
RUN apk update && apk upgrade && apk add --no-cache git

RUN mkdir -p /go/src/bitbucket.org/slonoed/chntr/
ADD . /go/src/bitbucket.org/slonoed/chntr
WORKDIR /go/src/bitbucket.org/slonoed/chntr

RUN go build
RUN go install

EXPOSE 3000

CMD ["chntr"]
You have mail in /var/mail/slonoed
