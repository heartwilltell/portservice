# Builder layer
FROM golang:1.20-alpine AS BUILDER

ARG BRANCH
ARG COMMIT

RUN apk add --no-cache ca-certificates build-base git
RUN mkdir -p /app
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-X main.Branch=$BRANCH -X main.Commit=$COMMIT" -o app ./cmd

# Final layer
FROM alpine:latest

COPY --from=BUILDER /app/app app

ENTRYPOINT ["./app"]

# Final layer
FROM alpine:latest

RUN addgroup -S appuser && adduser -S -G appuser appuser

WORKDIR /home/appuser

COPY --from=BUILDER /app/app .

RUN chown appuser:appuser /home/appuser

USER appuser

ENTRYPOINT ["/home/appuser/app"]