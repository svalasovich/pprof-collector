# Example

## Build project

```bash
go build
```

## Collect data from many servers

Follows pattern: `./pprof-collector <seconds to collect> <urls>`

```bash
./pprof-collector 30 http://localhost:8080/debug/pprof/profile http://localhost:8081/debug/pprof/profile http://localhost:8082/debug/pprof/profile
```

## View data in browser

Starts web server with required port

```bash
go tool pprof -http=:8092 result.pprof
```