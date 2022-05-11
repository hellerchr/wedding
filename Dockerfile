# build
FROM golang:alpine as build
RUN apk add --no-cache ca-certificates git
WORKDIR /hochzeit
COPY go.mod go.sum ./
RUN go mod download
ARG CACHEBUST=CACHEBUST_ARG
COPY src/ ./src
WORKDIR /hochzeit/src
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags '-extldflags "-static"' -o app

# image
FROM scratch
WORKDIR /app
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
ARG CACHEBUST=CACHEBUST_ARG
COPY --from=build /hochzeit/src/static/ ./static
COPY --from=build /hochzeit/src/app .
EXPOSE 8080
CMD ["./app"]