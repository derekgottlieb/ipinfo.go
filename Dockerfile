FROM golang:alpine as builder
WORKDIR /src/app
COPY . /src/app/
RUN go build -o ipinfo

FROM alpine
WORKDIR /app/
RUN apk add --no-cache wget & \
    wget https://github.com/P3TERX/GeoLite.mmdb/raw/download/GeoLite2-City.mmdb
COPY --from=builder /src/app /app
CMD ["./ipinfo"]
