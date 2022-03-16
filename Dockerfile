# Build container
FROM golang:1.17.8-buster AS build_base

RUN apt-get install git

WORKDIR /tmp/app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

# Unit tests
RUN CGO_ENABLED=0 go test -v

RUN go build -o ./out/trading-platform .

# Start fresh from a smaller image for the runtime container
FROM debian:buster

RUN apt-get update \
    && apt-get install -y --no-install-recommends sqlite3 ca-certificates

RUN update-ca-certificates

WORKDIR /app

# Copy the go executable from the build stage
COPY --from=build_base /tmp/app/out/trading-platform  /app/trading-platform 

EXPOSE ${HTTP_PORT}

CMD ./trading-platform  start-server