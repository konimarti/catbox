package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"git.sr.ht/~sircmpwn/getopt"
	"github.com/emersion/go-mbox"
)

func main() {
	opts, _, err := getopt.Getopts(os.Args, "c:")
	if err != nil {
		panic(err)
	}

	var filter string = "cat"
	for _, opt := range opts {
		if opt.Option == 'c' {
			filter = opt.Value
		}
	}

	r := mbox.NewReader(os.Stdin)
	for i := 1; ; i++ {

		// get reader for next message
		mr, err := r.NextMessage()
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			panic(err)
		}

		// create command for every message
		cmd, err := createCmd(filter)
		if err != nil {
			panic(err)
		}

		// update cmd's env
		env := os.Environ()
		env = append(env, fmt.Sprintf("NR=%06d", i))
		cmd.Env = env

		// set new reader as stdin to cmd
		cmdStdin, err := cmd.StdinPipe()
		if err != nil {
			panic(err)
		}
		go func() {
			defer cmdStdin.Close()
			io.Copy(cmdStdin, mr)
		}()

		// run command and collect output
		output, err := cmd.CombinedOutput()
		if err != nil {
			panic(err)
		}

		// flush command output to stdout
		io.Copy(os.Stdout, bytes.NewReader(output))
	}

}

func createCmd(s string) (*exec.Cmd, error) {
	s = strings.ReplaceAll(s, "\\$", "$")
	args := []string{"sh", "-c", s}

	cmd := exec.Command(args[0], args[1:]...)

	return cmd, nil
}
