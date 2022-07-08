package bf

import (
	"io"
)

// cell size could be 8 bits, 16, 32 or 64
// Make it stack-based.
// Make it a library.

// You must not parse the code before executing it. Important: the code must be read from a reader interface,
// such as io.Reader or io.ByteReader.
// Instead, the interpreter should read and execute one command at a time.
func Run(src io.Reader) error {
	insts, err := io.ReadAll(src)
	if err != nil {
		return err
	}

	return run(insts)
}

func run(insts []byte) error {
	s := &runner{}
	s.Init(insts)
	return s.run()
}

// Create a command-line executable brainfuck interpreter program.
// The program must be using your brainfuck library under the hood.

// ### Pluses
// * Automated tests.
// * Good documentation of your code.
// * Good error handling: errors that help the users of your library to understand what is wrong with their brainfuck code and where.
