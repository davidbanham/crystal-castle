package destiny

import (
	"crypto/md5"
	"fmt"
	"strconv"
)

const TARGET_HASHES = 2500

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

type Replicant struct {
	Hostname       string `json:"hostname"`
	HashAdjustment int    `json:"HashAdjustment"`
	Hash           string `json:"Hash"`
}

type Manifest map[string][]Job

type Plan map[string][]Job

type JobList []Job

type AdjustedJobList []Job

type NodeList []Node

type ReplicatedNodeList []Replicant

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

func AdjustJobs(jobs JobList, plan Plan, nodes ReplicatedNodeList) (ret AdjustedJobList) {

	adjusted := map[string]Job{}

	for node, jobList := range plan {
		for _, job := range jobList {
			job = adjustJob(job, node, nodes)
			adjusted[job.Job] = job
		}
	}

	for _, job := range jobs {
		ret = append(ret, adjusted[job.Job])
	}
	return ret
}

func adjustJob(job Job, target string, nodes ReplicatedNodeList) (ret Job) {
	ret = job
	winner := findMatchingNode(job, nodes)
	if winner == target {
		return ret
	}
	ret.HashAdjustment++
	return adjustJob(ret, target, nodes)
}

func findMatchingNode(job Job, nodes ReplicatedNodeList) (winningNode string) {
	highest := 0

	for _, node := range nodes {
		hash := hashString(node.Hash + job.Job + strconv.Itoa(job.HashAdjustment))

		sum := HashValue(hash)

		if sum > highest {
			highest = sum
			winningNode = node.Hostname
		}
	}
	return winningNode
}

func HashValue(data string) (int) {
	sum := 0
	for _, r := range []rune(data) {
		sum += int(r)
	}
	return sum
}

func BuildManifest(jobs AdjustedJobList, nodes NodeList, circle ReplicatedNodeList) (manifest Manifest) {
	manifest = make(Manifest)

	for _, node := range nodes {
		manifest[node.Hostname] = []Job{}
	}

	for _, job := range jobs {
		matchingNode := findMatchingNode(job, circle)
		manifest[matchingNode] = append(manifest[matchingNode], job)
	}

	return manifest
}

func BuildReplicatedNodeList(nodes NodeList) (ret ReplicatedNodeList) {
	for _, node := range nodes {
		num_replicas := TARGET_HASHES / len(nodes)

		i := num_replicas
		for i > 0 {
			newReplicant := Replicant{
				Hostname: node.Hostname,
				HashAdjustment: i,
				Hash: hashString(node.Hostname + strconv.Itoa(i)),
			}
			ret = append(ret, newReplicant)
			i--
		}
	}
	return ret
}

func hashString(in string) (string) {
	hash := md5.Sum([]byte(in))
	return fmt.Sprintf("%x", hash)
}
