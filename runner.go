package bf

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"unsafe"
)

var (
	ErrNoCommands       = errors.New("no commands to run")
	ErrMaxArraySize     = fmt.Errorf("backing array size can't exceeds %d", size)
	ErrNegativeIndex    = errors.New("array index can't be less than zero")
	ErrIllegalCharNul   = errors.New("illegal character NUL")
	ErrLoopDoesNotMatch = errors.New("loop openings/closing ([/]) count does not match")
)

const (
	eof  = -1 // end of file
	size = 30000
)

type offset struct {
	offset     int // character offset
	rdOffset   int // reading offset (position after current character)
	lineOffset int // current line offset
}

// runner represents the interpreter of Brainfuck. It scans each command and runs it right away.
// Since Go uses rune (int32) to represent character values, bf uses int32 for cell size, with a
// 30K item backing array to provided better error .
//
type runner struct {
	src []byte // source

	c    int         // backing array limit counter
	arr  [size]int32 // backing array
	ptr  unsafe.Pointer
	lpos []offset // loop stack positions

	// scanning state
	ch rune // current character
	o  offset
}

func (r *runner) Init(src []byte) error {
	r.src = src

	r.ch = ' '
	r.o.offset = 0
	r.o.rdOffset = 0
	r.o.lineOffset = 0
	r.ptr = unsafe.Pointer(&r.arr)

	err := validate(r.src)
	if err != nil {
		return err
	}

	o, err := r.next()
	if err != nil {
		return nextErr(err, o)
	}

	return nil
}

func (r *runner) next() (offset, error) {
	if r.o.rdOffset < len(r.src) {
		r.o.offset = r.o.rdOffset
		if r.ch == '\n' {
			r.o.lineOffset = r.o.offset
		}
		nr, w := rune(r.src[r.o.rdOffset]), 1
		if nr == 0 {
			return r.o, ErrIllegalCharNul
		}
		r.o.rdOffset += w
		r.ch = nr
	} else {
		r.o.offset = len(r.src)
		if r.ch == '\n' {
			r.o.lineOffset = r.o.offset
		}
		r.ch = eof
	}

	return offset{}, nil
}

func (r *runner) skipWhitespace() (offset, error) {
	for r.ch == ' ' || r.ch == '\t' || r.ch == '\n' || r.ch == '\r' {
		if o, err := r.next(); err != nil {
			return o, err
		}
	}

	return offset{}, nil
}

// exec reads and executes each command until it reaches EOF
func (r *runner) exec(out io.Writer) error {
	var res bytes.Buffer
	for {
		o, err := r.skipWhitespace()
		if err != nil {
			return nextErr(err, o)
		}

		// determine token value
		ch := r.ch
		o, err = r.next()
		if err != nil {
			return nextErr(err, o)
		}

		switch ch {
		case '-':
			*(*int32)(r.ptr)--
		case '[':
			// always push the position to the stack
			r.lpos = append([]offset{
				{
					offset:     r.o.offset - 2,
					rdOffset:   r.o.rdOffset - 2,
					lineOffset: r.o.lineOffset,
				},
			}, r.lpos...)
		case ']':
			if len(r.lpos) == 0 || (r.lpos[0].offset == 0 && r.lpos[0].rdOffset == 0 && r.lpos[0].lineOffset == 0) {
				return ErrLoopDoesNotMatch
			}

			// continue the loop of *ptr is not zero
			if *(*int32)(r.ptr) != 0 {
				// jump back, and read again
				r.o = r.lpos[0]
				o, err := r.next()
				if err != nil {
					return nextErr(err, o)
				}
			}

			// always pop the loop stack
			r.lpos = r.lpos[1:]
		case '>':
			r.c++
			if r.c >= size {
				return ErrMaxArraySize
			}
			r.ptr = unsafe.Add(r.ptr, unsafe.Sizeof(r.arr[0]))
		case '<':
			r.c--
			if r.c < 0 {
				return ErrNegativeIndex
			}
			r.ptr = unsafe.Pointer(uintptr(r.ptr) - unsafe.Sizeof(r.arr[0]))
		case '+':
			*(*int32)(r.ptr)++
		case '.':
			res.WriteString(string(*(*int32)(r.ptr)))

			// ignore unknown commands
		}

		if ch == eof {
			break
		}
	}

	_, err := out.Write(res.Bytes())
	if err != nil {
		return fmt.Errorf("failed writing to the output: %w", err)
	}

	return nil
}

func validate(src []byte) error {
	s := string(src)

	if s == "" {
		return ErrNoCommands
	}

	var c int
	for _, r := range s {
		if r == '[' {
			c++
			continue
		}

		if r == ']' {
			c--
		}
	}

	if c != 0 {
		return ErrLoopDoesNotMatch
	}

	return nil
}

func nextErr(err error, o offset) error {
	return fmt.Errorf("%v at %d:%d (line:offset)", err, o.lineOffset, o.offset)
}
