FROM golang:latest

WORKDIR /home/golang

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN go build -v -o /usr/local/bin/app ./...

CMD ["app"]