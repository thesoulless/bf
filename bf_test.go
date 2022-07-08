package bf

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"testing"
	"unsafe"

	qt "github.com/frankban/quicktest"
)

func TestBF_Run(t *testing.T) {
	c := qt.New(t)

	t.Run("from string", func(t *testing.T) {
		s := `++++++++[>++++[>++>+++>+++>+<<
<<-]>+>+>->>+[<]<-]>>.>---.+++++++..+++.>
>.<-.<.+++.------.--------.>>+.>++.`
		input := strings.NewReader(s)
		want := "Hello World!\n"

		var buf []byte
		out := bytes.NewBuffer(buf)

		bfi, err := New(input, out, nil)
		c.Assert(err, qt.IsNil)

		err = bfi.Exec()

		c.Assert(err, qt.IsNil)
		c.Assert(out.String(), qt.Equals, want)
	})

	t.Run("negative index", func(t *testing.T) {
		s := `++++++++<<<<`
		input := strings.NewReader(s)
		wantErr := ErrNegativeIndex

		var buf []byte
		out := bytes.NewBuffer(buf)

		bfi, err := New(input, out, nil)
		c.Assert(err, qt.IsNil)
		err = bfi.Exec()

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

		_, err := New(input, out, nil)
		c.Assert(err, qt.IsNotNil)
		c.Assert(out.String(), qt.Equals, "")
		c.Assert(err, qt.Equals, wantErr)
	})

	t.Run("inputs", func(t *testing.T) {
		s := `,++       Cell c0 = 5
> ,++  Cell c1 = 6
[<+>-]<.`
		input := strings.NewReader(s)
		want := []byte{11}

		inps := "3\n4\n"
		args := strings.NewReader(inps)

		var buf []byte
		out := bytes.NewBuffer(buf)

		bfi, err := New(input, out, args)
		c.Assert(err, qt.IsNil)

		c.Assert(err, qt.IsNil)
		err = bfi.Exec()

		c.Assert(err, qt.IsNil)
		c.Assert(out.Bytes(), qt.ContentEquals, want)
	})
}

func TestBF_AddCommand(t *testing.T) {
	c := qt.New(t)

	t.Run("custom command", func(t *testing.T) {
		s := `++       Cell c0 = 2
> ++  Cell c1 = 2
[<+>-]<^.`
		input := strings.NewReader(s)
		want := []byte{16}

		var buf []byte
		out := bytes.NewBuffer(buf)

		bfi, err := New(input, out, nil)
		c.Assert(err, qt.IsNil)

		err = bfi.AddCommand('^', func(ptr unsafe.Pointer) {
			*(*int32)(ptr) *= *(*int32)(ptr)
		})
		c.Assert(err, qt.IsNil)
		err = bfi.Exec()

		c.Assert(err, qt.IsNil)
		c.Assert(out.Bytes(), qt.ContentEquals, want)
	})

	t.Run("duplicate command", func(t *testing.T) {
		s := `,++       Cell c0 = 5
> ,++  Cell c1 = 6
[<+>-]<.`
		input := strings.NewReader(s)

		var buf []byte
		out := bytes.NewBuffer(buf)

		bfi, err := New(input, out, nil)
		c.Assert(err, qt.IsNil)
		err = bfi.AddCommand('^', func(ptr unsafe.Pointer) {
			*(*int32)(ptr) *= *(*int32)(ptr)
		})
		c.Assert(err, qt.IsNil)
		err = bfi.AddCommand('^', func(ptr unsafe.Pointer) {
			*(*int32)(ptr)++
		})
		c.Assert(err, qt.IsNotNil)
	})
}

func TestBF_RemoveCommand(t *testing.T) {
	c := qt.New(t)
	t.Run("remove command", func(t *testing.T) {
		s := `++       Cell c0 = 2
> ++  Cell c1 = 2
[<+>-]<.`
		input := strings.NewReader(s)

		var buf []byte
		out := bytes.NewBuffer(buf)

		bfi, err := New(input, out, nil)
		c.Assert(err, qt.IsNil)
		err = bfi.AddCommand('^', func(ptr unsafe.Pointer) {
			*(*int32)(ptr) *= *(*int32)(ptr)
		})
		c.Assert(err, qt.IsNil)
		bfi.RemoveCommand('^')
		err = bfi.AddCommand('^', func(ptr unsafe.Pointer) {
			*(*int32)(ptr)++
		})
		c.Assert(err, qt.IsNil)
	})
}

func Example() {
	s := `,+++++[>++++[>++>+++>+++>+<<
<<-]>+>+>->>+[<]<-]>>.>---.+++++++..+++.>
>.<-.<.+++.------.--------.>>+.>++.`
	input := strings.NewReader(s)

	var buf []byte
	out := bytes.NewBuffer(buf)

	inps := "3\n"
	args := strings.NewReader(inps)

	bfi, err := New(input, out, args)
	if err != nil {
		log.Printf("error: %v", err)
		return
	}
	err = bfi.Exec()
	if err != nil {
		log.Printf("error: %v", err)
		return
	}
	fmt.Println(out.String())

	// Output: Hello World!
}
