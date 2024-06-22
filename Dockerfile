FROM golang:alpine3.20 as build
WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
ARG ENV="dev"
RUN go build -v -o ./bin/api -ldflags "-X env.EXT_ENVIRONMENT=${ENV}" ./cmd/main.go

FROM debian:latest
COPY --from=build /usr/src/app/bin/api /bin/api
EXPOSE 3000
CMD ["/bin/api"]

