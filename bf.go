// Package bf is a library that reads a series of Brainfuck commands from
// an io.Reader and writes the outputs of those commands to an io.Writer
package bf

import (
	"io"
)

// Run reads the Brainfuck commands from src and writes the outputs to out
func Run(src io.Reader, out io.Writer) error {
	insts, err := io.ReadAll(src)
	if err != nil {
		return err
	}

	return run(insts, out)
}

func run(insts []byte, out io.Writer) error {
	s := &runner{}
	err := s.Init(insts)
	if err != nil {
		return err
	}
	return s.exec(out)
}
