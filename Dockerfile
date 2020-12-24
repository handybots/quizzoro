FROM golang:alpine as builder

WORKDIR /src

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o quizzoro

FROM alpine

WORKDIR /app

COPY --from=builder /src/quizzoro .

ENTRYPOINT ["/app/quizzoro"]
