FROM golang:1.27-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /geopulse ./cmd/

FROM alpine:3.20
RUN apk add --no-cache ca-certificates
COPY --from=build /geopulse /geopulse
EXPOSE 8080
CMD ["/geopulse"]
