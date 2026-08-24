FROM golang:1.22-alpine AS build
WORKDIR /src
COPY . .
RUN go build -o /industrial-server ./cmd/server
FROM alpine:3.20
COPY --from=build /industrial-server /industrial-server
EXPOSE 8080
ENTRYPOINT ["/industrial-server"]
