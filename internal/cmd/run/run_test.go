package run

import (
	"io/ioutil"
	"os"
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestCmd(t *testing.T) {
	c := qt.New(t)

	t.Run("from string", func(t *testing.T) {
		input := "++++++++[>++++[>++>+++>+++>+<<<<-]>+>+>->>+[<]<-]>>.>---.+++++++..+++.>>.<-.<.+++.------.--------.>>+.>++."
		want := "Hello World!\n"

		cmd := Cmd()
		err := cmd.Flags().Set("string", input)
		cmd.Args = nil
		if err != nil {
			t.Fatalf("failed to set flags: %v", err)
		}

		oStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		err = cmd.Execute()

		w.Close()
		out, _ := ioutil.ReadAll(r)
		os.Stdout = oStdout

		c.Assert(err, qt.IsNil)
		c.Assert(string(out), qt.Equals, want)
	})

	t.Run("from file", func(t *testing.T) {
		file := "../../../testdata/test.bf"
		want := "Hello World!\n"

		cmd := Cmd()
		err := cmd.Flags().Set("file", file)
		cmd.Args = nil
		if err != nil {
			t.Fatalf("failed to set flags: %v", err)
		}

		oStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		err = cmd.Execute()

		w.Close()
		out, _ := ioutil.ReadAll(r)
		os.Stdout = oStdout

		c.Assert(err, qt.IsNil)
		c.Assert(string(out), qt.Equals, want)
	})
}
