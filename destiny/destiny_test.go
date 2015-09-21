package destiny

import (
	"encoding/json"
	"log"
	"os"
	"reflect"
	"testing"
)

var jobs JobList
var nodes NodeList

var manifest, rebalanced Manifest
var target Plan

var adjusted AdjustedJobList

func logErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func TestMain(m *testing.M) {
	_jobs, err := os.Open("fixtures/jobs.json")
	logErr(err)

	_nodes, err := os.Open("fixtures/nodes.json")
	logErr(err)

	_manifest, err := os.Open("fixtures/manifest.json")
	logErr(err)

	_target, err := os.Open("fixtures/target.json")
	logErr(err)

	_rebalanced, err := os.Open("fixtures/rebalanced.json")
	logErr(err)

	_adjusted, err := os.Open("fixtures/adjusted.json")
	logErr(err)

	jsonParser := json.NewDecoder(_jobs)
	jsonParser.Decode(&jobs)

	jsonParser = json.NewDecoder(_nodes)
	jsonParser.Decode(&nodes)

	jsonParser = json.NewDecoder(_manifest)
	jsonParser.Decode(&manifest)

	jsonParser = json.NewDecoder(_target)
	jsonParser.Decode(&target)

	jsonParser = json.NewDecoder(_rebalanced)
	jsonParser.Decode(&rebalanced)

	jsonParser = json.NewDecoder(_adjusted)
	jsonParser.Decode(&adjusted)

	retCode := m.Run()
	os.Exit(retCode)
}

func TestJson(t *testing.T) {
	t.Skip()
	t.Log(jobs)
	t.Log(nodes)
	t.Log(manifest)
	t.Log(target)
	t.Log(rebalanced)
	t.Log(adjusted)
}

func TestBuildPlan(t *testing.T) {
	plan := BuildPlan(jobs, nodes)
	equal := reflect.DeepEqual(plan, target)
	if !equal {
		t.Fail()
	}
}

func TestBuildReplicatedNodeList(t *testing.T) {
	replicated := BuildReplicatedNodeList(nodes)
	numReplicants := len(replicated)
	if (numReplicants != TARGET_HASHES) {
		t.Fail()
	}
}

func TestHashValue(t *testing.T) {
	hash := "59bcc3ad6775562f845953cf01624225"
	targetSum := 2059
	result := HashValue(hash)
	if (result != targetSum) { t.Fail() }
}

func TestAdjustJob(t *testing.T) {
	testJob := Job{
		Job: "Threatening Oryx",
		HashAdjustment: 0,
	}
	replicatedNodes := BuildReplicatedNodeList(nodes)
	adjusted := adjustJob(testJob, "one", replicatedNodes)
	target := Job{
		Job: "Threatening Oryx",
		HashAdjustment: 4,
	}
	equal := reflect.DeepEqual(adjusted, target)
	if !equal {
		t.Fail()
	}
}

func TestAdjustJobs(t *testing.T) {
	plan := BuildPlan(jobs, nodes)
	replicatedNodes := BuildReplicatedNodeList(nodes)
	adjustedJobs := AdjustJobs(jobs, plan, replicatedNodes)
	equal := reflect.DeepEqual(adjustedJobs, adjusted)
	if !equal {
		t.Fail()
	}
}

func TestFindMatchingNode(t *testing.T) {
	replicatedNodes := BuildReplicatedNodeList(nodes)
	testJob := Job{
		Job: "Threatening Oryx",
		HashAdjustment: 4,
	}
	winner := findMatchingNode(testJob, replicatedNodes)
	if winner != "one" {
		t.Fail()
	}
}

func TestBuildManifest(t *testing.T) {
	plan := BuildPlan(jobs, nodes)
	replicatedNodes := BuildReplicatedNodeList(nodes)
	adjustedJobs := AdjustJobs(jobs, plan, replicatedNodes)

	builtManifest := BuildManifest(adjustedJobs, nodes, replicatedNodes)
	equal := reflect.DeepEqual(builtManifest, manifest)
	if !equal {
		t.Fail()
	}
}
