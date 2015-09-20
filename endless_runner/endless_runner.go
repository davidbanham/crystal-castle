package endless_runner

import (
	"os"
	"os/exec"
)

func Run(name string, args ...string) (chan *os.Process, chan error, chan bool) {

	//Make a communication channel that transmits processes
	commChannel := make(chan *os.Process)

	// Make a channel to pass back errors on
	errChannel := make(chan error)

	// Make a channel to listen for stop events
	stopChannel := make(chan bool)

	go waiter(commChannel, errChannel, stopChannel, name, args...)

	return commChannel, errChannel, stopChannel
}

func runProc(name string, args ...string) (*exec.Cmd, error) {
	cmd := exec.Command(name, args...)

	//cmd.Stdout = os.Stdout
	//cmd.Stderr = os.Stderr

	err := cmd.Start()
	return cmd, err
}

// Waiter takes a channel that speaks processes, a command name, and some arguments
// It runs that command with the arguments and passes teh resulting process back across the channel
// It them waits for that process to exit. When it does, it recurses. This causes the process to be restarted and the new process passed back along the channel
func waiter(comm chan *os.Process, errs chan error, stops chan bool, name string, args ...string) {
	cmd, err := runProc(name, args...)

	shouldStop := false

	go func() {
		shouldStop = <-stops
		cmd.Process.Kill()
		cmd.Process.Release()
		stops <- true
	}()

	if err != nil {
		errs <- err
		return
	} else {
		comm <- cmd.Process
		cmd.Wait()
		if shouldStop {
			return
		}
		go waiter(comm, errs, stops, name, args...)
	}
}
