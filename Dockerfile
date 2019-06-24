FROM golang:1.12 as serverbuilder
WORKDIR /loki
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main

FROM scratch

COPY --from=serverbuilder /loki/main  /server/main

WORKDIR  /server

CMD ["./main"]
