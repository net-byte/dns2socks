FROM golang:alpine

WORKDIR /app
COPY . /app
ENV GO111MODULE=on
RUN go build -o ./bin/dns2socks ./main.go

ENTRYPOINT ["./bin/dns2socks"]

