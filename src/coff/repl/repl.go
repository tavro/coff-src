package repl

import (
	"bufio"
	"fmt"
	"io"
	"coff-src/src/coff/lexer"
	"coff-src/src/coff/parser"
	"coff-src/src/coff/eval"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	
	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		
		if !scanned {
			return
		}
		
		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)
		program := p.ParseProgram()
	
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}
	
		evaluated := eval.Eval(program)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}