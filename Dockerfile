FROM golang:1.22-alpine AS build
WORKDIR /src
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /risk-engine ./cmd/server
FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /risk-engine /risk-engine
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/risk-engine"]
