// Package bf is a library that reads a series of Brainfuck commands from
// an io.Reader and writes the outputs of those commands to an io.Writer
package bf

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
	"unsafe"
)

var (
	size       = 3
	defaultCms = []rune{'>', '<', '+', '-', '.', ',', '[', ']'}

	ErrNoCommands       = errors.New("no commands to run")
	ErrDuplicateCmd     = errors.New("duplicate command")
	ErrNegativeIndex    = errors.New("array index can't be less than zero")
	ErrIllegalCharNul   = errors.New("illegal character NUL")
	ErrLoopDoesNotMatch = errors.New("loop openings/closing ([/]) count does not match")
)

const (
	eof = -1 // end of file
)

type offset struct {
	offset     int // character offset
	rdOffset   int // reading offset (position after current character)
	lineOffset int // current line offset
}

// BF represents the interpreter of Brainfuck. It scans each command and runs it right away.
// Since Go uses rune (int32) to represent character values, bf uses int32 for cell size, starting
// with a 30K item backing array. The backing array will expand in case of a need so the bf be
// Turing complete.
//
// On validation, it returns error on empty command set, and when loop beginning and endings
// does not match.
// On execution, it returns error on moving index to negative, and encountering a NUL character.
type BF struct {
	src []byte // source
	out io.Writer

	c     int     // backing array limit counter
	arr   []int32 // backing array
	ptr   unsafe.Pointer
	lpos  []offset // loop stack positions
	ucmds map[rune]func(unsafe.Pointer)

	// scanning state
	ch      rune // current character
	o       offset
	inp     io.Reader // input (,) reader
	inpscan *bufio.Scanner
}

// New creates a new BF. It returns error on reading from src, or validating commands.
// It takes src as the source of commands, out as where to write the outputs, and
// input as where it should read the inputs (, command)
func New(src io.Reader, out io.Writer, input io.Reader) (*BF, error) {
	insts, err := io.ReadAll(src)
	if err != nil {
		return nil, err
	}

	s := &BF{src: insts, out: out, inp: input}
	err = s.init()
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (b *BF) init() error {
	b.ch = ' '
	b.arr = make([]int32, size)
	b.o.offset = 0
	b.o.rdOffset = 0
	b.o.lineOffset = 0
	b.ptr = unsafe.Pointer(&b.arr[0])
	b.inpscan = bufio.NewScanner(b.inp)
	b.ucmds = make(map[rune]func(unsafe.Pointer))

	err := validate(b.src)
	if err != nil {
		return err
	}

	o, err := b.next()
	if err != nil {
		return nextErr(err, o)
	}

	return nil
}

// AddCommand associate a function to a character. It returns error on overriding
// current valid commands (defaults and user-defined). It passes the pointer to the
// current array cell to the function.
func (b *BF) AddCommand(cmd rune, f func(ptr unsafe.Pointer)) error {
	for _, c := range defaultCms {
		if cmd == c {
			return ErrDuplicateCmd
		}
	}
	if _, ok := b.ucmds[cmd]; ok {
		return ErrDuplicateCmd
	}

	b.ucmds[cmd] = f
	return nil
}

func (b *BF) next() (offset, error) {
	if b.o.rdOffset < len(b.src) {
		b.o.offset = b.o.rdOffset
		if b.ch == '\n' {
			b.o.lineOffset = b.o.offset
		}
		nr, w := rune(b.src[b.o.rdOffset]), 1
		if nr == 0 {
			return b.o, ErrIllegalCharNul
		}
		b.o.rdOffset += w
		b.ch = nr
	} else {
		b.o.offset = len(b.src)
		if b.ch == '\n' {
			b.o.lineOffset = b.o.offset
		}
		b.ch = eof
	}

	return offset{}, nil
}

func (b *BF) skipWhitespace() (offset, error) {
	for b.ch == ' ' || b.ch == '\t' || b.ch == '\n' || b.ch == '\r' {
		if o, err := b.next(); err != nil {
			return o, err
		}
	}

	return offset{}, nil
}

// Exec reads and executes each command until it reaches EOF
func (b *BF) Exec() error {
	var res bytes.Buffer
	for {
		o, err := b.skipWhitespace()
		if err != nil {
			return nextErr(err, o)
		}

		// determine token value
		ch := b.ch
		o, err = b.next()
		if err != nil {
			return nextErr(err, o)
		}

		switch ch {
		case '-':
			*(*int32)(b.ptr)--
		case '[':
			// always push the position to the stack
			b.lpos = append([]offset{
				{
					offset:     b.o.offset - 2,
					rdOffset:   b.o.rdOffset - 2,
					lineOffset: b.o.lineOffset,
				},
			}, b.lpos...)
		case ']':
			if len(b.lpos) == 0 || (b.lpos[0].offset == 0 && b.lpos[0].rdOffset == 0 && b.lpos[0].lineOffset == 0) {
				return ErrLoopDoesNotMatch
			}

			// continue the loop of *ptr is not zero
			if *(*int32)(b.ptr) != 0 {
				// jump back, and read again
				b.o = b.lpos[0]
				o, err := b.next()
				if err != nil {
					return nextErr(err, o)
				}
			}

			// always pop the loop stack
			b.lpos = b.lpos[1:]
		case '>':
			b.c++
			if b.c >= size {
				// instead of returning error expand the backing array cap
				size = len(b.arr) + size
				b.arr = append(make([]int32, 0, size), b.arr...)
			}
			b.ptr = unsafe.Add(b.ptr, unsafe.Sizeof(b.arr[0]))
		case '<':
			b.c--
			if b.c < 0 {
				return ErrNegativeIndex
			}
			b.ptr = unsafe.Pointer(uintptr(b.ptr) - unsafe.Sizeof(b.arr[0]))
		case '+':
			*(*int32)(b.ptr)++
		case ',':
			fmt.Print("Enter value (int32): ")
			b.inpscan.Scan()
			text := b.inpscan.Text()
			//reader := bufio.NewReader(b.inp)
			//text, _ := reader.ReadString('\n')
			if text != "" {
				i, err := strconv.ParseInt(text, 10, 32)
				if err != nil {
					return fmt.Errorf("invalid input: %v", err)
				}
				*(*int32)(b.ptr) = int32(i)
			}
		case '.':
			res.WriteString(string(*(*int32)(b.ptr)))

			// ignore unknown commands
		}

		if f, ok := b.ucmds[ch]; ok {
			f(b.ptr)
		}

		if ch == eof {
			break
		}
	}

	_, err := b.out.Write(res.Bytes())
	if err != nil {
		return fmt.Errorf("failed writing to the output: %v", err)
	}

	return nil
}

// validate the commands source. It returns error on empty command set, and when loop
// beginning and endings does not match.
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
