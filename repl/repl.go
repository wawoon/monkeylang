package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/wawoon/monkeylang/lexer"
	"github.com/wawoon/monkeylang/token"
)

const PROMPT = "> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			break
		}

		line := scanner.Text()
		l := lexer.New(line)

		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Printf("%s\n", tok.String())
		}
	}
}
