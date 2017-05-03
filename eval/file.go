package eval

import (
	"bufio"
	"bytes"
	"os"
)

type FileObject struct {
	File    *os.File
	Name    string
	Scanner *bufio.Scanner
}

func (f *FileObject) Inspect() string  { return f.Name }
func (f *FileObject) Type() ObjectType { return FILE_OBJ }
func (f *FileObject) CallMethod(method string, args ...Object) Object {
	switch method {
	case "close":
		f.File.Close()
		return NULL
	case "read":
		return f.Read(args...)
	case "readline":
		return f.ReadLine()
	default:
		return newError(NOMETHODERROR, method, f.Type())
	}
}

func (f *FileObject) Read(args ...Object) Object {
	if len(args) != 0 {
		return newError(ARGUMENTERROR, "0", len(args))
	}
	fs := bufio.NewScanner(f.File)
	var out bytes.Buffer
	for {
		scanned := fs.Scan()
		out.WriteString(fs.Text())
		if !scanned {
			break
		}
		if err := fs.Err(); err != nil {
			return &Error{Message: err.Error()}
		}
	}
	return &String{Value: out.String()}
}

func (f *FileObject) ReadLine(args ...Object) Object {
	if len(args) != 0 {
		return newError(ARGUMENTERROR, "0", len(args))
	}
	if f.Scanner == nil {
		f.Scanner = bufio.NewScanner(f.File)
		f.Scanner.Split(bufio.ScanLines)
	}
	line := f.Scanner.Scan()
	if err := f.Scanner.Err(); err != nil {
		return &Error{Message: err.Error()}
	}
	if !line {
		return NULL
	}
	return &String{Value: f.Scanner.Text()}
}
