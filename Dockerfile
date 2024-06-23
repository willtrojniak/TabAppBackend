FROM golang:alpine3.20 as build_base
WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .

FROM build_base as build_migrate
ARG ENV="dev"
RUN go build -v -o ./bin/migrate -ldflags "-X env.EXT_ENVIRONMENT=${ENV}" ./cmd/migrate/main.go

FROM build_base as build
ARG ENV="dev"
RUN go build -v -o ./bin/api -ldflags "-X env.EXT_ENVIRONMENT=${ENV}" ./cmd/main.go

FROM debian:bookworm
RUN apt update
RUN apt install ca-certificates -y
RUN update-ca-certificates
COPY --from=build_migrate /usr/src/app/bin/migrate /bin/migrate
COPY --from=build /usr/src/app/bin/api /bin/api
EXPOSE 3000
CMD ["/bin/api"]

