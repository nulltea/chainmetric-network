FROM golang:1.16-alpine AS build

WORKDIR /go/src/github.com/timoth-y/chainmetric-contracts/
COPY ../../.. .

RUN go mod vendor
RUN go build -o service -v ./src/identity

FROM alpine:3.11 as prod

COPY --from=build /go/src/github.com/timoth-y/chainmetric-contracts/service /app/service
COPY --from=build /go/src/github.com/timoth-y/chainmetric-contracts/src/identity/data  /app/src/identity/data

USER 1000

WORKDIR /app
CMD ./service
