FROM golang:1.16-alpine AS build

WORKDIR /go/src/github.com/timoth-y/chainmetric-contracts/
COPY ../../.. .

RUN go mod vendor
RUN go build -o service -v ./src/identity

FROM alpine:3.11 as prod

COPY --from=build /go/src/github.com/timoth-y/chainmetric-contracts/identity /app/service

USER 1000

WORKDIR /app
CMD ./service
