package bf

import (
	"errors"
	"fmt"
	"unsafe"
)

const (
	eof = -1 // end of file
)

type offset struct {
	offset     int // character offset
	rdOffset   int // reading offset (position after current character)
	lineOffset int // current line offset
}

type runner struct {
	src []byte // source

	arr  [30000]int32
	ptr  unsafe.Pointer
	lpos []offset // loop stack positions

	// scanning state
	ch rune // current character
	o  offset
}

func (r *runner) Init(src []byte) {
	r.src = src

	r.ch = ' '
	r.o.offset = 0
	r.o.rdOffset = 0
	r.o.lineOffset = 0
	r.ptr = unsafe.Pointer(&r.arr)

	r.next()
}

func (r *runner) next() {
	if r.o.rdOffset < len(r.src) {
		r.o.offset = r.o.rdOffset
		if r.ch == '\n' {
			r.o.lineOffset = r.o.offset
		}
		nr, w := rune(r.src[r.o.rdOffset]), 1
		if nr == 0 {
			panic("illegal character NUL")
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
}

func (r *runner) skipWhitespace() {
	for r.ch == ' ' || r.ch == '\t' || r.ch == '\n' || r.ch == '\r' {
		r.next()
	}
}

func (r *runner) run() error {
	for {
		r.skipWhitespace()

		// determine token value
		ch := r.ch
		r.next()

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
				return errors.New("broken loop")
			}

			// continue the loop of *ptr is not zero
			if *(*int32)(r.ptr) != 0 {
				// jump back, and read again
				r.o.offset = r.lpos[0].offset
				r.o.rdOffset = r.lpos[0].rdOffset
				r.o.lineOffset = r.lpos[0].lineOffset
				r.next()
			}

			// always pop the loop stack
			r.lpos = r.lpos[1:]
		case '>':
			r.ptr = unsafe.Add(r.ptr, unsafe.Sizeof(r.arr[0]))
		case '<':
			r.ptr = unsafe.Pointer(uintptr(r.ptr) - unsafe.Sizeof(r.arr[0]))
		case '+':
			*(*int32)(r.ptr)++
		case '.':
			fmt.Printf("%s", string(*(*int32)(r.ptr)))

			// ignore unknown commands
		}

		if ch == eof {
			break
		}
	}

	return nil
}
