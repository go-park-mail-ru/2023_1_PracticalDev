FROM golang:1.19.6-alpine3.17 AS base

WORKDIR /src
ENV CGO_ENABLED=0
COPY go.* ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

FROM base AS build
RUN --mount=target=. \
    --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -o /out/search cmd/search/*.go

FROM golang:1.19.6-alpine3.17
COPY --from=build /out/search /
ENTRYPOINT "/search"
