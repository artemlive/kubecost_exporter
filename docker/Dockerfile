FROM golang:1.17-alpine AS build

RUN apk update && apk add curl \
                          git \
                          bash \
                          make

WORKDIR /app

COPY . .

RUN make build

# binary itself
FROM golang:1.17-alpine as deploy

RUN apk add bash && addgroup -S nonroot && adduser -s /bin/bash nonroot -G nonroot -D

ENV KUBECOST_URL="http://kubecost.example.com"

WORKDIR /app

COPY --from=build /app/bin/kubecost_exporter /usr/local/bin/

ENTRYPOINT ["/usr/local/bin/kubecost_exporter"]
EXPOSE 9150

USER nonroot:nonroot
