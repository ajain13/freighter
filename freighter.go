/**
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/paypal/gorealis"
)

type JobJson struct {
	NAME      string  `json:"name"`
	CPU       float64 `json:"cpu"`
	RAM       int64   `json:"ram"`
	IMAGE     string  `json:"image"`
	COMMAND   string  `json:"cmd"`
	INSTANCES int32   `json:"inst"`
}

func (j *JobJson) Validate() bool {
	fmt.Println(*j)
	if j.CPU <= 0.0 {
		return false
	}
	if j.RAM <= 0 {
		return false
	}
	if j.INSTANCES <= 0 {
		return false
	}
	return true
}

func main() {

	url := flag.String("url", "", "URL at which the Aurora Scheduler exists as [url]:[port]")
	jsonFile := flag.String("file", "", "JSON file containing job definition")
	username := flag.String("username", "aurora", "Username to use for authorization")
	password := flag.String("password", "secret", "Password to use for authorization")
	flag.Parse()

	if *jsonFile == "" {
		flag.Usage()
		os.Exit(1)
	}

	file, err := ioutil.ReadFile(*jsonFile)

	if err != nil {
		fmt.Println("Error opening file", err)
		os.Exit(1)
	}

	var jsonJob []JobJson
	json.Unmarshal(file, &jsonJob)

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

	r, err := realis.NewRealisClient(
		realis.BasicAuth(*username, *password),
		realis.ThriftBinary(),
		realis.SchedulerUrl(*url))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer r.Close()

	var auroraJob realis.Job

	for _, job := range jsonJob {
		auroraJob = realis.NewJob().
			Environment("prod").
			Role("benchmarks").
			Name(job.NAME).
			CPU(job.CPU).
			RAM(job.RAM).
			Disk(400).
			IsService(false).
			InstanceCount(jsonJob[0].INSTANCES).
			AddPorts(1)

		fmt.Println("Creating docker based job : ", job.NAME)
		container := realis.NewDockerContainer().Image(job.IMAGE).AddParameter("network", "host")
		auroraJob.Container(container)
		resp, err := r.CreateJob(auroraJob)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println(resp.String())
	}
}
