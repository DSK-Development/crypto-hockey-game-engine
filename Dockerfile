FROM golang:1.26-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /out/engine ./cmd/engine
# cache-bust: 2026-06-02

FROM gcr.io/distroless/static:nonroot
COPY --from=build /out/engine /engine
EXPOSE 8081
USER nonroot:nonroot
ENTRYPOINT ["/engine"]
