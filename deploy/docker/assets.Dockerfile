FROM golang:1.16-alpine AS build

WORKDIR /go/src/github.com/timoth-y/chainmetric-network/
COPY ../.. .

RUN go build -o chaincode -v ./src/assets

FROM alpine:3.11 as prod

COPY --from=build /go/src/github.com/timoth-y/chainmetric-network/chaincode /app/chaincode

USER 1000

WORKDIR /app
CMD ./chaincode
