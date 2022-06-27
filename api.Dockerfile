# Image
FROM golang:alpine as builder

# Set workdir
WORKDIR /app

# Copy over files
COPY . .

# Install dependencies
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o main ./cmd/api/main.go

# Minimize using busybox
FROM busybox

WORKDIR /app

COPY --from=builder /app/main /usr/bin/
COPY --from=builder /app/.env /app

# Add script to wait for MYSQL to finish first before booting (crucial)
COPY ./entrypoint.sh /entrypoint.sh
ADD https://raw.githubusercontent.com/eficode/wait-for/v2.1.0/wait-for /usr/local/bin/wait-for
RUN chmod +rx /usr/local/bin/wait-for /entrypoint.sh

ENTRYPOINT [ "sh", "/entrypoint.sh" ]
