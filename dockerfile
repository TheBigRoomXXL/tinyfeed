# ╔═════════════════════════════════════════════════╗
# ║                  BUILD STAGE                    ║
# ╚═════════════════════════════════════════════════╝
FROM golang:1.23.1-alpine3.20 AS build

# Install the dependencies
RUN apk add upx
COPY go.mod go.sum ./
RUN go mod download

# Build static binary
COPY *.go ./
COPY built-in ./
RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o /tinyfeed
RUN upx /tinyfeed

# Create the user file 
RUN adduser -H -D feed
# USER feed

# ╔═════════════════════════════════════════════════╗
# ║               PRODUCTION STAGE                  ║
# ╚═════════════════════════════════════════════════╝
FROM scratch AS production

# Create a user to avoid runing as root
COPY --from=build /etc/passwd /etc/passwd
USER feed

# Where the input and output files are
WORKDIR /app

# Copy the dependencies
COPY --chown=feed:feed --from=build /tinyfeed /usr/local/bin/tinyfeed

# Copy the certificates authorities
COPY --from=alpine:latest /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["/usr/local/bin/tinyfeed"]
