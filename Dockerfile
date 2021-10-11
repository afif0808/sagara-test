FROM golang:alpine

RUN apk update && apk add --no-cache git

WORKDIR /media/afif0808/data/goprojects/github.com/afif0808/sagara-test

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

EXPOSE 8080

