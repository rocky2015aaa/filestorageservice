FROM golang:1.22 AS build

RUN mkdir /app
ADD . /app
WORKDIR /app

# Define build-time arguments
ARG VERSION=dev
ARG BUILD=dev
ARG DATE=1970-01-01_00:00:00

# Build the Go application with build-time arguments
RUN CGO_ENABLED=0 go build -ldflags "-X github.com/rocky2015aaa/filestorageservice/internal/config.Version=${VERSION} -X github.com/rocky2015aaa/filestorageservice/internal/config.Build=${BUILD} -X github.com/rocky2015aaa/filestorageservice/internal/config.Date=${DATE}" -o filestorageservice cmd/filestorageservice/main.go

FROM alpine:latest

RUN mkdir /application
WORKDIR /application

COPY --from=build /app/filestorageservice .
COPY --from=build /app/docs /application/docs
COPY --from=build /app/.env .

EXPOSE 8080

CMD ["./filestorageservice"]
