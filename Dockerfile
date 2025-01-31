FROM golang:1.23 AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /main

FROM scratch
COPY --from=builder /main /main
ENV DOCKER_DB_HOST="host.docker.internal"
ENV DB_PORT="5432"
ENV DB_USER="root"
ENV DB_PASS="root"
ENV DB_NAME="root"
ENV GIN_MODE="release"
ENTRYPOINT [ "/main" ]
