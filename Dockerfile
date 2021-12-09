FROM golang:1.17 as build

WORKDIR /app

COPY go.sum ./
COPY go.mod ./

RUN go get ./...

ADD . .

RUN CGO_ENABLED=0 go build -o build/program
# COPY settings.docker.yaml ./build

FROM alpine

WORKDIR /app

COPY --from=build /app/build/* ./

# ENTRYPOINT ./program serve --config ${APP_CONFIG_FILE}
ENTRYPOINT ./program serve
