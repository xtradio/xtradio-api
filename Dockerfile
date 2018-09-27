FROM golang:alpine AS build
WORKDIR /src
COPY xtradio-api.go .

RUN apk update && apk add git ca-certificates

RUN go get -d -v

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/xtradio-api .

FROM scratch

LABEL "maintainer"="XTRadio Ops <contact@xtradio.org"
LABEL "version"="0.1"
LABEL "description"="XTRadio API"

COPY --from=build /src/bin/xtradio-api /bin/xtradio-api

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# ADD ./bin/xtradio-api /xtradio-api

# #ADD ./bin/xtradio-api /xtradio-api
EXPOSE 10000

CMD ["/xtradio-api"]
