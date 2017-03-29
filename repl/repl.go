package repl

import (
	"io"
	"monkey/lexer"
	"monkey/parser"
	"os"
	"path/filepath"

	"github.com/peterh/liner"
)

const PROMPT = ">> "

func Start(out io.Writer) {
	history := filepath.Join(os.TempDir(), ".monkey_history")
	l := liner.NewLiner()
	defer l.Close()

	l.SetCtrlCAborts(true)

	if f, err := os.Open(history); err == nil {
		l.ReadHistory(f)
		f.Close()
	}

	for {
		if line, err := l.Prompt(PROMPT); err == nil {
			if line == "exit" {
				if f, err := os.Create(history); err == nil {
					l.WriteHistory(f)
					f.Close()
				}
				break
			}
			lex := lexer.New(line)
			p := parser.New(lex)

			program := p.ParseProgram()
			if len(p.Errors()) != 0 {
				printParserErrors(out, p.Errors())
				continue
			}

			io.WriteString(out, program.String())
			io.WriteString(out, "\n")
			l.AppendHistory(line)
		}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
