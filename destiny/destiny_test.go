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
}

func TestBuildPlan(t *testing.T) {
	plan := BuildPlan(jobs, nodes)
	equal := reflect.DeepEqual(plan, target)
	if !equal {
		t.Fail()
	}
}
