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
	args := make([]string, len(os.Args))
	copy(args, os.Args)
	opts, optind, err := getopt.Getopts(args, "c:h")
	if err != nil {
		panic(err)
	}

	var filter string = "cat"
	for _, opt := range opts {
		if opt.Option == 'c' {
			filter = opt.Value
		}
		if opt.Option == 'h' {
			usage()
			return
		}
	}

	var inputReader io.ReadCloser = os.Stdin
	if len(args[optind:]) > 0 {
		name := strings.Join(args[optind:], " ")
		if _, err := os.Stat(name); err == nil {
			inputReader, err = os.Open(name)
			if err != nil {
				panic(err)
			}
			defer inputReader.Close()
		}

	}

	r := mbox.NewReader(inputReader)
	for i := 1; ; i++ {

		// get reader for next message
		mr, err := r.NextMessage()
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			handleErr(i, err)
			continue
		}

		// create command for every message
		cmd, err := createCmd(filter)
		if err != nil {
			handleErr(i, err)
			continue
		}

		// update cmd's env
		cmd.Env = append(os.Environ(), fmt.Sprintf("NR=%06d", i))

		// set new reader as stdin to cmd
		cmdStdin, err := cmd.StdinPipe()
		if err != nil {
			handleErr(i, err)
			continue
		}
		go func() {
			defer cmdStdin.Close()
			io.Copy(cmdStdin, mr)
		}()

		// run command and collect output
		output, err := cmd.CombinedOutput()
		if err != nil {
			handleErr(i, err)
			continue
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

func handleErr(i int, err error) {
	fmt.Printf("skipping messgae %d. Error encountered: %v", i, err)
}

func usage() {
	usage := `
	Usage: catbox [-h|-c <cmd>] <mbox>

	Options:
		-h	Print usage.
		-c cmd	Specify shell command.
	
	catbox will pipe every message from an mbox file as an input to a shell
	command. A message counter is available as a shell variable $NR.
	
	If no file is specified, catbox will read from stdin.
	`
	fmt.Println(usage)
}
