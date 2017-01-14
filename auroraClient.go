package main
import (
        "github.com/rdelval/gorealis"
        "encoding/json"
        "flag"
        "fmt"
        "os"
        "io/ioutil"
)

type JobJson struct {
        NAME      string            `json:"name"`
        CPU       float64           `json:"cpu"`
        RAM       int64             `json:"ram"`
        WATTS     float64           `json:"watts"`
	IMAGE     string            `json:"image"`
        COMMAND   string            `json:"cmd"`
        INSTANCES int32             `json:"inst"`
}

func (j *JobJson) Validate() bool {
  if j.CPU <= 0.0 {
    return false
  }
  if j.RAM <= 0 {
    return false
  }
  if j.WATTS <= 0.0 {
    return false
  }
  if j.INSTANCES <= 0 {
    return false
  }
  return true
}


func main() {

  url := flag.String("url", "", "URL at which the Aurora Scheduler exists as [url]:[port]")
  jsonFile := flag.String("file","","JSON file containing job definition")
  clustersConfig := flag.String("clusters", "", "Location of the clusters.json file used by aurora.")
  clusterName := flag.String("cluster", "devcluster", "Name of cluster to run job on (only necessary if clusters is set)")
  username := flag.String("username", "aurora", "Username to use for authorization")
  password := flag.String("password", "secret", "Password to use for authorization")
  flag.Parse()

  if *clustersConfig != "" {
                clusters, err := realis.LoadClusters(*clustersConfig)
                if err != nil {
                        fmt.Println(err)
                        os.Exit(1)
                }

                cluster, ok := clusters[*clusterName]
                if !ok {
                        fmt.Printf("Cluster %s chosen doesn't exist\n", *clusterName)
                        os.Exit(1)
                }

                *url, err = realis.LeaderFromZK(cluster)
                if err != nil {
                        fmt.Println(err)
                        os.Exit(1)
                }
  }


  if *jsonFile == "" {
    flag.Usage()
    os.Exit(1)
  }

  file,err :=ioutil.ReadFile(*jsonFile)

  if err != nil {
    fmt.Println("Error opening file",err)
    os.Exit(1)
  }

  var jsonJob []JobJson
  json.Unmarshal(file,&jsonJob)

  if err != nil {
    fmt.Println("Error parsing file ", err)
    os.Exit(1)
  }

  for _, job := range jsonJob {
    err := job.Validate()
    if !err {
      fmt.Println("Invalid job !")
      os.Exit(1)
    }
  }

  config, err := realis.NewDefaultConfig(*url)
  if err != nil {
        fmt.Println(err)
        os.Exit(1)
  }

  if *username != "" && *password != "" {

	realis.AddBasicAuth(&config, *username, *password)
  }

  r := realis.NewClient(config)
  defer r.Close()

  var aurora_job realis.Job

  for _,job := range jsonJob {
		aurora_job = realis.NewJob().
	                  Environment("prod").
	                  Role("benchmarks").
	                  Name(job.NAME).
	                  CPU(job.CPU).
	                  RAM(job.RAM).
	                  Disk(400).
	                  IsService(false).
	                  InstanceCount(jsonJob[0].INSTANCES).
	                  AddPorts(1)

	  fmt.Println("Creating docker based job : ",job.NAME)
	  container := realis.NewDockerContainer().Image(job.IMAGE).AddParameter("network","host")
	  aurora_job.Container(container)
	  resp, err := r.CreateJob(aurora_job)
	  if err !=nil {
	      fmt.Println(err)
	      os.Exit(1)
	  }
	  fmt.Println(resp.String())


  }



}
