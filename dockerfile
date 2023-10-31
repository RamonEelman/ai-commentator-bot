# Use an official Go runtime as a parent image
FROM golang:1.21 AS build-env

# Set working directory inside the container
WORKDIR /app

# Download dependencies
COPY dca/go.mod /app/dca/
COPY fakeyou/go.mod /app/fakeyou/
COPY go.mod go.sum /app/
RUN go mod download

# Copy local code to the container
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o myapp .

### Final stage
FROM alpine:latest
WORKDIR /app
COPY --from=build-env /app/myapp /app/
CMD ["./myapp"]
