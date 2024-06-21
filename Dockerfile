FROM golang:alpine3.20 as build
WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go build -v -o ./bin/api ./cmd/main.go

FROM debian:latest
COPY --from=build /usr/src/app/bin/api /bin/api
EXPOSE 8080
CMD ["/bin/api"]

