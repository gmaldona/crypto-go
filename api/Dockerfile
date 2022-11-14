FROM golang:latest AS build

WORKDIR /home/golang

COPY . .

RUN go mod download && go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -a -o app .

FROM alpine:latest

RUN apk update && \
    apk upgrade && \
    apk add ca-certificates

WORKDIR /app/

COPY --from=build /home/golang/app ./

RUN env && pwd && find .

CMD ["./app"]