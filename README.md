# freighter
Client for [Apache Aurora](https://github.com/apache/aurora) using [gorealis](https://github.com/paypal/gorealis) to launch docker container jobs in bulk.

## Using the client
```
Usage of ./auroraClient:
  -url string
        URL at which the Aurora Scheduler exists as [url]:[port]
  -file string
        JSON file containing job definition
  -username string
        Username to use for authorization (default "aurora")
  -password string
        Password to use for authorization (default "secret")
```

## Sample Command
```
go run freighter.go -file sample_freighter_workload.json -url http://192.168.33.7:8081
```

## Requirements

Aurora Scheduler must have the following options enabled for this tool to work:
```
  -allow_docker_parameters=true
  -require_docker_use_executor=false
```
