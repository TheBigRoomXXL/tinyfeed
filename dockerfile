# ╔═════════════════════════════════════════════════╗
# ║                  BUILD STAGE                    ║
# ╚═════════════════════════════════════════════════╝
FROM golang:1.23.1-alpine3.20 AS build

# Install the dependencies
COPY go.mod go.sum ./
RUN go mod download

# Build the binary
COPY *.go ./
COPY built-in ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /tinyfeed


# ╔═════════════════════════════════════════════════╗
# ║               PRODUCTION STAGE                  ║
# ╚═════════════════════════════════════════════════╝
FROM alpine:3.20 AS production

# Create a user to avoid runing as root
RUN adduser -H -D feed
USER feed

# Where the input and output files are
WORKDIR /app

# Copy the dependencies
COPY --chown=feed:feed --from=build /tinyfeed /usr/local/bin/tinyfeed

ENTRYPOINT ["/usr/local/bin/tinyfeed"]
