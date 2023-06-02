FROM golang AS build_base

WORKDIR /app
COPY . .
ENV GO111MODULE=on
ENV GOFLAGS=-mod=vendor
RUN go mod vendor && go build -o /app/app .


# FROM alpine:3.9 
FROM ubuntu
ENV GIN_MODE=release
WORKDIR /app
COPY --from=build_base /app/app /app/app

EXPOSE 8080

ENTRYPOINT ["/app/app"]