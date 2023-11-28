package client

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sandbox/scheduler"
	"time"

	"gonum.org/v1/gonum/stat/distuv"
)

type Client struct {
	Name         string
	Cluster      scheduler.Cluster
	SchedulerURL string
	maxJobMem    uint
	maxJobCore   uint
	time_dist    string
	job_dist     string
}

func Run(URL string) {
	client.SchedulerURL = URL
	client.newClient()
	client.sendJobs()
}

func (c *Client) newClient() {
	// request cluster information & this can include more information
	// in the future
	// key exchange ?
	res, err := http.Get(c.SchedulerURL + "/newClient")
	if err != nil {
		log.Println(err)
		return
	}
	var cluster scheduler.Cluster
	dec := json.NewDecoder(res.Body)
	err = dec.Decode(&cluster)
	if err != nil {
		log.Println(err)
		return
	}
	c.Cluster = cluster
	fmt.Printf("cluster %+v \n", cluster)
	// TODO: change this to user input
	c.time_dist = "poisson"

	// set job range
	c.setMaxCluster()
}

func (c *Client) setMaxCluster() {
	var maxCores uint
	var maxMem uint

	for _, node := range c.Cluster.Nodes {
		if node.Memory > maxMem {
			maxMem = node.Memory
		}
		if len(node.Cores) > int(maxCores) {
			maxCores = uint(len(node.Cores))
		}
	}

	c.maxJobCore = maxCores
	c.maxJobMem = maxMem
}

func (c *Client) sendJobs() {
	// approximating normal distribution
	job_dist := distuv.Beta{
		Alpha: 2,
		Beta:  2,
	}
	var id uint
	getJob := func() scheduler.Job {

		var j scheduler.Job
		j.Id = id
		fmt.Print(c.maxJobCore)
		j.CoresNeeded = uint(job_dist.Rand()) * c.maxJobCore // uint(rand.Intn(int(c.maxJobCore))) //uint( dist.Rand() * float64(c.maxJobCore))
		j.Duration = 10
		j.MemoryNeeded = uint(job_dist.Rand()) * c.maxJobMem // uint(rand.Intn(int(c.maxJobMem)))

		id++
		return j
	}

	// create time distrribution
	switch c.time_dist {
	case "poisson":
		time_dist := distuv.Poisson{
			Lambda: 10,
		}

		for {
			// each for loop is one minute, lambda jobs per minute is sent
			jobs := int(time_dist.Rand())
			fmt.Printf("jobs per minute: %v", jobs)
			time_between_jobs := 60 / jobs
			i := 0
			for i < jobs {
				i++
				// uint(dist.Rand() * float64(c.maxJobMem))
				//fmt.Printf("cores: %v memory %v\n", dist.Rand(), dist.Rand())
				j := getJob()
				SendJob(j)
				time.Sleep(time.Duration(time_between_jobs) * time.Second)

			}

		}

		/*
			i := 0
			for i < 10 {
				//var dist distuv.Normal

				dist := distuv.LogNormal{
					Mu:    1, // Mean of the normal distribution
					Sigma: 1, // Standard deviation of the normal distribution
				}

				var j scheduler.Job
				j.Id = uint(i)
				fmt.Print(c.maxJobCore)
				j.CoresNeeded = uint(rand.Intn(int(c.maxJobCore))) //uint( dist.Rand() * float64(c.maxJobCore))
				j.Duration = 10
				j.MemoryNeeded = uint(rand.Intn(int(c.maxJobMem))) // uint(dist.Rand() * float64(c.maxJobMem))
				fmt.Printf("cores: %v memory %v\n", dist.Rand(), dist.Rand())
				SendJob(j)
				i++
			}
		*/

		/*
			case "weibull":
				time_dist := distuv.Weibull{
					Lambda: 10,
					K:      3,
				}

				for {
					// each for loop is one minute, lambda jobs per minute is sent
					jobs := int(time_dist.Rand())
					fmt.Printf("jobs per minute: %v", jobs)
					time_between_jobs := 60 / jobs
					i := 0
					for i < jobs {
						i++
						// uint(dist.Rand() * float64(c.maxJobMem))
						//fmt.Printf("cores: %v memory %v\n", dist.Rand(), dist.Rand())
						j := getJob()
						SendJob(j)
						time.Sleep(time.Duration(time_between_jobs) * time.Second)

					}
				}
			}
		*/

	}
}
