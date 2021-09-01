FROM golang:1.16-alpine AS build

WORKDIR /go/src/github.com/timoth-y/chainmetric-contracts/
COPY ../.. .


RUN go build -o chaincode -v ./src/readings

FROM alpine:3.11 as prod

COPY --from=build /go/src/github.com/timoth-y/chainmetric-contracts/chaincode /app/chaincode

USER 1000

WORKDIR /app
CMD ./chaincode
