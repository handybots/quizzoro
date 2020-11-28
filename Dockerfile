FROM golang:1.15.5-alpine3.12

WORKDIR /tmp/build
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build -o /usr/bin/quizzoro .

WORKDIR /usr/bin
ENTRYPOINT ["quizzoro"]