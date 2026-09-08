# syntax=docker/dockerfile:1
FROM golang:1.27.1-bookworm AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
ARG VERSION=dev
ARG COMMIT=unknown
ARG BUILD_TIME=unknown
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath \
    -ldflags "-s -w -X main.version=${VERSION} -X main.commit=${COMMIT} -X main.buildTime=${BUILD_TIME}" \
    -o /out/goblog ./ && \
    CGO_ENABLED=0 GOOS=linux go build -trimpath -o /out/goblog-admin ./cmd/admin

FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /app
COPY --from=build /out/goblog /out/goblog-admin /app/bin/
COPY --from=build /src/templates /app/templates
COPY --from=build /src/static /app/static
COPY --from=build /src/migrations /app/migrations
COPY --from=build /src/blog.sql /src/LICENSE /src/README.md /src/SECURITY.md /app/
COPY --from=build /src/conf/config.example.yml /app/conf/config.yml
ENV GOBLOG_CONFIG=/app/conf/config.yml
EXPOSE 9091
ENTRYPOINT ["/app/bin/goblog"]
