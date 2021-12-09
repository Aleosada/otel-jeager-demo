# OpenTelemetry Jeager Demo

## Run instructions

1. Run docker containers

`docker-compose up -d --build`

2. Run client

`go run . client -r {numberOfRequests} -i {interval}`

3. Open the url http://localhost:16686 (Jeager UI) and check the results
