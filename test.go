package tranquil

import (
	"fmt"
	"io/ioutil"
	"path"
	"runtime"
	"strings"
	"testing"
)

var (
	prevDesc string
	reset    = "\033[0m"
	white    = "\033[37m\033[1m"
	grey     = "\x1B[90m"
	red      = "\033[31m\033[1m"
)

type errorInfo struct {
	lines    []string
	filename string
	number   int
}

// Test represents a behavior test (e.g. an 'it').
type Test struct {
	Description string
	T           *testing.T
	ExpectFn    func(Expect)
	Setup       *Setup
}

// Run executes the behavior test.
func (t *Test) Run() {
	t.Setup.execBeforeEachFns()
	t.ExpectFn(func(val interface{}) *Assertion {
		return NewAssertion(t, val)
	})
	t.Setup.execAfterEachFns()
}

// PrintError prints an error if a test fails.
func (t *Test) PrintError(message string) {
	if prevDesc != t.Description {
		fmt.Printf("%s    %s\n", white, t.Description)
		prevDesc = t.Description
	}

	errorInfo, err := t.getErrorInfo()
	if err != nil {
		return
	}

	fmt.Printf("%s    %s %s %s %s\n", red, message, grey, path.Base(errorInfo.filename), reset)
	fmt.Printf("%s        %d. %s\n", grey, errorInfo.number-1, errorInfo.lines[0])
	fmt.Printf("%s        %d. %s %s\n", white, errorInfo.number, errorInfo.lines[1], reset)
	fmt.Printf("%s        %d. %s\n", grey, errorInfo.number+1, errorInfo.lines[2])
	fmt.Println(reset)
	t.T.Fail()
}

func (*Test) getErrorInfo() (errorInfo, error) {
	file := ""
	line := 0

	var ok bool
	for i := 0; ; i++ {
		_, file, line, ok = runtime.Caller(i)
		if !ok {
			return errorInfo{}, fmt.Errorf("")
		}
		if strings.HasSuffix(file, "_test.go") {
			break
		}
	}

	f, err := ioutil.ReadFile(file)
	if err != nil {
		return errorInfo{}, fmt.Errorf("")
	}

	lines := strings.Split(string(f), "\n")[line-2 : line+2]
	for i, l := range lines {
		lines[i] = strings.Replace(l, "\t", "  ", -1)
	}

	returnVal := errorInfo{
		lines:    lines,
		filename: file,
		number:   int(line),
	}
	return returnVal, nil
}
