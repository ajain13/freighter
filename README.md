# freighter
Client for [Apache Aurora](https://github.com/apache/aurora) using [gorealis](https://github.com/paypal/gorealis) to launch docker container jobs in bulk.

## Using the client
```
Usage of ./auroraClient:
  -url string
        URL at which the Aurora Scheduler exists as [url]:[port]
  -file string
        JSON file containing job definition
  -clusters string
        Location of the clusters.json file used by aurora.
  -cluster string
        Name of cluster to run job on (default "devcluster")
  -username string
        Username to use for authorization (default "aurora")
  -password string
        Password to use for authorization (default "secret")
```

## Sample Command
```
go run auroraClient.go -file sample_electron_workload.json -url http://192.168.33.7:8081
```
