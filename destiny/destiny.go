package destiny

import (
)

type Job struct {
	Job            string `json:"job"`
	HashAdjustment int    `json:"hashAdjustment"`
}

type Node struct {
	Ip struct {
		External string `json:"external"`
		Internal string `json:"internal"`
	}
	Hostname string `json:"hostname"`
}

type Manifest map[string][]Job

type Plan map[string][]Job

type JobList []Job

type AdjustedJobList []Job

type NodeList []Node

func BuildPlan(jobs JobList, nodes NodeList) (plan Plan) {
	plan = make(Plan)

	i := 0
	for _, job := range jobs {
		node := nodes[i].Hostname

		if plan[node] == nil {
			plan[node] = []Job{}
		}

		plan[node] = append(plan[node], job)

		i++
		if i > len(nodes)-1 {
			i = 0
		}
	}
	return plan
}

func AdjustJobs(jobs JobList, plan Plan) (ret AdjustedJobList) {
	return ret
}

func BuildManifest(jobs AdjustedJobList, node NodeList) (manifest Manifest) {
	return manifest
}
