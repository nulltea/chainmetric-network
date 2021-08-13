# This image is a microservice in golang for the Degree chaincode
FROM golang:1.16-alpine AS build

WORKDIR /go/src/github.com/timoth-y/chainmetric-contracts/
COPY ../../.. .

# Build application
RUN go build -o service -v ./src/users

# Production ready image
FROM alpine:3.11 as prod

COPY --from=build /go/src/github.com/timoth-y/chainmetric-contracts/chaincode /app/service

USER 1000

WORKDIR /app
CMD ./service
