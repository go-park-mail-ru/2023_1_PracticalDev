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
    go build -o /out/pickpin cmd/pickpin/*.go

FROM base AS unit-test
ENV PGHOST=localhost
ENV REDIS_HOST=localhost
RUN --mount=target=. \
    --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    mkdir /out && go test -coverprofile=/out/cover.out ./...

FROM scratch AS coverage
COPY --from=unit-test /out/cover.out /cover.out

FROM golang:1.19.6-alpine3.17 AS print-coverage
WORKDIR /src
COPY --from=unit-test /out/cover.out /cover.out
RUN --mount=target=.,rw=true \
    --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    /bin/sh -c ' \
    cat /cover.out | fgrep -v "easyjson" > /cover.out.filtered && \
    go tool cover -func=/cover.out.filtered > /out.stat \
    '
CMD [ "cat", "/out.stat" ]

FROM golangci/golangci-lint:v1.51.2-alpine AS lint-base

FROM base AS lint
RUN --mount=target=. \
    --mount=from=lint-base,src=/usr/bin/golangci-lint,target=/usr/bin/golangci-lint \
    --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/root/.cache/golangci-lint \
    golangci-lint run --timeout 10m0s ./...

FROM golang:1.19.6-alpine3.17
COPY --from=build /out/pickpin /
ENTRYPOINT "/pickpin"
