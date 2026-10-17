curl localhost:8080/health
curl "localhost:8080/hello?name=gopher"
curl localhost:8080/panic
curl -X POST "localhost:8080/loglevel?level=debug"
curl localhost:8080/panic
curl -X POST localhost:8080/shutdown
