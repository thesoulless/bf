package bf

import (
	"bytes"
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestRun(t *testing.T) {
	c := qt.New(t)

	t.Run("from string", func(t *testing.T) {
		s := `++++++++[>++++[>++>+++>+++>+<<
<<-]>+>+>->>+[<]<-]>>.>---.+++++++..+++.>
>.<-.<.+++.------.--------.>>+.>++.`
		input := strings.NewReader(s)
		want := "Hello World!\n"

		var buf []byte
		out := bytes.NewBuffer(buf)

		err := Run(input, out)

		c.Assert(err, qt.IsNil)
		c.Assert(out.String(), qt.Equals, want)
	})

	t.Run("negative index", func(t *testing.T) {
		s := `++++++++<<<<`
		input := strings.NewReader(s)
		wantErr := ErrNegativeIndex

		var buf []byte
		out := bytes.NewBuffer(buf)

		err := Run(input, out)

		c.Assert(err, qt.IsNotNil)
		c.Assert(out.String(), qt.Equals, "")
		c.Assert(err, qt.Equals, wantErr)
	})

	t.Run("invalid loops", func(t *testing.T) {
		s := `++>+++++[[->+<]`
		input := strings.NewReader(s)
		wantErr := ErrLoopDoesNotMatch

		var buf []byte
		out := bytes.NewBuffer(buf)

		err := Run(input, out)

		c.Assert(err, qt.IsNotNil)
		c.Assert(out.String(), qt.Equals, "")
		c.Assert(err, qt.Equals, wantErr)
	})
}
