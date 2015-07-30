package endless_runner

import (
  "os/exec"
  "os"
)

func Run(name string, args ...string) (chan *os.Process) {

  //Make a communication channel that transmits processes
  commChannel := make(chan *os.Process)

  go waiter(commChannel, name, args...)

  return commChannel
}

func runProc(name string, args ...string) *exec.Cmd {
  cmd := exec.Command(name, args...)

  //cmd.Stdout = os.Stdout
  //cmd.Stderr = os.Stderr

  err := cmd.Start()
  if err != nil {
    panic(err)
  }
  return cmd
}

// Waiter takes a channel that speaks processes, a command name, and some arguments
// It runs that command with the arguments and passes teh resulting process back across the channel
// It them waits for that process to exit. When it does, it recurses. This causes the process to be restarted and the new process passed back along the channel
func waiter(comm chan *os.Process, name string, args ...string) {
  cmd := runProc(name, args...)
  comm <- cmd.Process
  cmd.Wait()
  go waiter(comm, name, args...)
}
