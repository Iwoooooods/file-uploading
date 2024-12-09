# Build stage
FROM golang:1.23-alpine AS build

WORKDIR /server

# Copy the source code into the image
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o cmd/server/main.go

# test stage
FROM build as test
RUN make test

# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /server /server

EXPOSE 8080

CMD ["make", "run.server", "PORT=8080"]