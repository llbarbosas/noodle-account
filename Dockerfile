FROM golang:1.16 as build-env

WORKDIR /tmp/app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o ./server cmd/start.go

FROM gcr.io/distroless/base

WORKDIR /app

COPY --from=build-env /tmp/app/server server
COPY --from=build-env /tmp/app/http/certs http/certs
COPY --from=build-env /tmp/app/http/public http/public
COPY --from=build-env /tmp/app/http/views http/views

EXPOSE 3001

CMD ["./server"]