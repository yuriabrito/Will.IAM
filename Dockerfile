FROM golang:1.11-alpine

RUN mkdir -p /app/bin

RUN apk update
RUN apk add postgresql git make musl-dev gcc

ADD . /go/src/github.com/ghostec/Will.IAM
RUN cd /go/src/github.com/ghostec/Will.IAM && \
  make setup-project && \
  make build && \
  mv bin/Will.IAM /app/Will.IAM && \
  mv config /app/config && \
  mv Makefile /app/Makefile
  
WORKDIR /app

EXPOSE 4040

CMD /app/Will.IAM start-api
