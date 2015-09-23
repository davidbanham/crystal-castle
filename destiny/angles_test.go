package destiny

import (
	"encoding/json"
	"io/ioutil"
	"strconv"
	"testing"
)

func TestWriteArcPlot(t *testing.T) {
	t.Skip()
	replicated := BuildReplicatedNodeList(nodes)
	angles := []int{}
	for _, replicant := range replicated {
		angles = append(angles, replicant.Angle)
	}
	encoded, _ := json.Marshal(angles)
	ioutil.WriteFile("./angles.json", encoded, 0644)
	jobAngles := []int{}
	for _, job := range adjusted {
		jobAngle := toAngle(HashValue(hashString(job.Job+strconv.Itoa(job.HashAdjustment))) * ANGLE_MULT_FACTOR)
		jobAngles = append(jobAngles, jobAngle)
	}
	encoded, _ = json.Marshal(jobAngles)
	ioutil.WriteFile("./jobAngles.json", encoded, 0644)
}
