FROM --platform=$BUILDPLATFORM golang:alpine3.20 AS build_base
WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .

FROM build_base AS build_migrate
ENV GOCACHE=/root/.cache/go-build
ARG ENV
RUN --mount=type=cache,target="/root/.cach/go-build" GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -v -o ./bin/migrate -ldflags "-X github.com/willtrojniak/TabAppBackend/env.EXT_ENVIRONMENT=${ENV}" ./cmd/migrate/main.go

FROM build_base AS build
ENV GOCACHE=/root/.cache/go-build
ARG ENV
RUN --mount=type=cache,target="/root/.cach/go-build" GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -v -o ./bin/api -ldflags "-X github.com/willtrojniak/TabAppBackend/env.EXT_ENVIRONMENT=${ENV}" ./cmd/main.go

FROM --platform=$BUILDPLATFORM debian:bookworm
RUN apt update
RUN apt install ca-certificates -y
RUN update-ca-certificates
COPY --from=build_migrate /usr/src/app/bin/migrate /bin/migrate
COPY --from=build /usr/src/app/bin/api /bin/api
EXPOSE 3000
CMD ["/bin/api"]

