FROM golang:1.11-alpine as serverbuilder
WORKDIR /go/src/github.com/cjburchell/restmock
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main

FROM scratch

COPY --from=serverbuilder /go/src/github.com/cjburchell/restmock/main  /server/main

WORKDIR  /server

CMD ["./main"]
