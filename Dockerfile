FROM golang:1.8-alpine

RUN apk update
RUN apk add nodejs=6.9.5-r1

ENV NODE_ENV=production
ENV GO_ENV=production

RUN mkdir -p /go/src/github.com/slonoed/chantra/
WORKDIR /go/src/github.com/slonoed/chantra
ADD ./package.json /go/src/github.com/slonoed/chantra/package.json
RUN npm i
ADD . /go/src/github.com/slonoed/chantra

RUN go build
RUN go install

RUN npm run build

EXPOSE 3000

CMD ["chantra"]
