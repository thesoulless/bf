package bf

import (
	"io"
)

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

// ### Pluses
// * Automated tests.
// * Good documentation of your code.
// * Good error handling: errors that help the users of your library to understand what is wrong with their brainfuck code and where.
