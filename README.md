# BF

BF is a [Brainfuck](https://en.wikipedia.org/wiki/Brainfuck) interpreter.

## Installation
`go get -u github.com/thesoulless/bf`

## How to use

```go
package main

import (
	"fmt"
	"log"
	"bytes"
	"strings"
	
	"github.com/thesoulless/bf"
)

func main() {
	s := `,+++++[>++++[>++>+++>+++>+<<
<<-]>+>+>->>+[<]<-]>>.>---.+++++++..+++.>
>.<-.<.+++.------.--------.>>+.>++.`
	input := strings.NewReader(s)

	var buf []byte
	out := bytes.NewBuffer(buf)

	inps := "3\n"
	args := strings.NewReader(inps)

	bfi, err := bf.New(input, out, args)
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
```

You can also use the binaries as follows:

`$ bf run -s "BF_COMMANDS_HERE"`

`$ bf run -f ./path/to/file.bf`
