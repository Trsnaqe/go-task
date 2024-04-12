# Build the application from source
FROM golang:1.22.0 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /api ./cmd/main.go

# Run the tests in the container
FROM build-stage AS run-test-stage
RUN go test -v ./...

# Deploy the application binary into a lean image
FROM scratch AS build-release-stage

WORKDIR /

COPY --from=build-stage /api /api

# Copy migrations into the container
COPY ./cmd/migrate/migrations /migrations


# Run migrations on startup
ENTRYPOINT ["/api"]

