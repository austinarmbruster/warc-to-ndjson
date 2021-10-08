# Inspiration from https://medium.com/@chemidy/create-the-smallest-and-secured-golang-docker-image-based-on-scratch-4752223b7324
FROM golang:1.17-alpine as builder

# Create appuser
RUN adduser -D -g '' appuser

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

#RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go test ./... && \
    #CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s"

RUN CGO_ENABLED=0 GOOS=linux go test ./... && \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o theapp

FROM scratch

# Import the user and group files from the builder.
COPY --from=builder /etc/passwd /etc/passwd

COPY --from=builder /app/theapp /app/theapp

USER appuser

ENTRYPOINT ["/app/theapp"]