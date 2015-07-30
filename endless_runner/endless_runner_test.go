package endless_runner

import (
  "testing"
  "os/exec"
  "os"
  "fmt"
  "bufio"
  "strings"
  "strconv"
)

func handleError(err error) {
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
}

func countProcs(target string) int {
  ps := exec.Command("ps")
  stdout, err := ps.StdoutPipe()
  handleError(err)

  err = ps.Start()
  handleError(err)

  defer ps.Wait()

  buff := bufio.NewScanner(stdout)
  var allText []string

  for buff.Scan() {
    allText = append(allText, buff.Text())
  }

  matchedCount := 0
  for _, item := range allText {
    matched := strings.Contains(item, target)
    if (matched) { matchedCount++ }
  }

  return matchedCount
}

// That that Run actually runs the given process
func TestRun(t *testing.T) {
  // Run a meaningless process that never exits
  procChan, _ := Run("tail", "-f", "/dev/null")
  // Get the process out of the communication channel
  process := <-procChan
  specificCount := countProcs(strconv.Itoa(process.Pid))
  if (specificCount != 1) {
    t.Error("Reported pid was not found running by ps")
  }
  _, procErr := os.FindProcess(process.Pid)
  if (procErr != nil) {
    t.Error("Process not found")
  }
  process.Kill()
}

// See what happens when you pass an invalid command
func TestFail(t *testing.T) {
  procChan, errChan := Run("asdasfafsa", "asdasf")
  select {
  case <-procChan:
    t.Error("Shouldn't have returned anything on the proc channel")
  case err := <-errChan:
    if err.Error() != "exec: \"asdasfafsa\": executable file not found in $PATH" {
      t.Error("Returned error doesn't look right")
    }
  }
}
