package bf

import (
	"github.com/thesoulless/bf/token"
	"log"
	"unsafe"
)

const (
	eof = -1 // end of file
)

type runner struct {
	src   []byte // source
	stack []token.Token

	insts [30000]byte
	ptr   unsafe.Pointer

	// scanning state
	ch         rune // current character
	offset     int  // character offset
	rdOffset   int  // reading offset (position after current character)
	lineOffset int  // current line offset
}

func (r *runner) Init(src []byte) {
	r.src = src

	r.ch = ' '
	r.offset = 0
	r.rdOffset = 0
	r.lineOffset = 0
	r.ptr = unsafe.Pointer(&r.insts)

	r.next()
}

func (r *runner) next() {
	if r.rdOffset < len(r.src) {
		r.offset = r.rdOffset
		if r.ch == '\n' {
			r.lineOffset = r.offset
		}
		nr, w := rune(r.src[r.rdOffset]), 1
		if nr == 0 {
			panic("illegal character NUL")
		}
		r.rdOffset += w
		r.ch = nr
	} else {
		r.offset = len(r.src)
		if r.ch == '\n' {
			r.lineOffset = r.offset
		}
		r.ch = eof
	}
}

func (r *runner) skipWhitespace() {
	for r.ch == ' ' || r.ch == '\t' || r.ch == '\n' || r.ch == '\r' {
		r.next()
	}
}

func (r *runner) peek() byte {
	if r.rdOffset < len(r.src) {
		return r.src[r.rdOffset]
	}
	return 0
}

func (r *runner) Scan() (tok token.Token) {
	r.skipWhitespace()

	// determine token value
	ch := r.ch

	r.next()
	switch ch {
	case -1:
		tok = token.EOF
	case '-':
		tok = token.SUB
	case '[':
		tok = token.LBRACK
	case ']':
		tok = token.RBRACK
	case '>':
		tok = token.GTR
	case '<':
		tok = token.LSS
	case '+':
		tok = token.ADD
	default:
		log.Printf("illegal character %#U", ch)
		tok = token.INVALID
	}

	return tok
}

func (r *runner) run() error {
	var (
		tok token.Token
	)

	e := true
	for {
		r.skipWhitespace()

		// determine token value
		ch := r.ch
		r.next()

		switch ch {
		case -1:
			tok = token.EOF
		case '-':
			tok = token.SUB
		case '[':
			tok = token.LBRACK
		case ']':
			tok = token.RBRACK
		case '>':
			tok = token.GTR
			r.stack = append([]token.Token{token.GTR}, r.stack...)
		case '<':
			tok = token.LSS
		case '+':
			tok = token.ADD
			r.stack = append([]token.Token{token.ADD}, r.stack...)
		default:
			log.Printf("illegal character %#U", ch)
			tok = token.INVALID
		}

		if e {
			for _, i := range r.stack {
				switch i {
				case token.ADD:
					*(*byte)(r.ptr)++
				case token.GTR:
					r.ptr = unsafe.Add(r.ptr, unsafe.Sizeof(r.insts[0]))
				}
			}
			r.stack = nil
		}

		if tok == token.EOF {
			break
		}
	}

	log.Printf("tok: %s", tok.String())

	log.Printf("0: %v", r.insts[0])
	log.Printf("0: %v", r.insts[1])
	return nil
}
