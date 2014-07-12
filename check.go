package gophtu

import (
	"bytes"
	"os"
	"runtime"
	"strings"
	"testing"
	"text/template"
)

var templ *template.Template

const (
	msgFmt = `{{if .ExpRes}}` +
		`want: {{printf "'%v'" .Expected}}, got: {{printf "'%v'" .Got}}{{else}}` +
		`want: result different than {{printf "'%v'" .Expected}}{{end}}` +
		`{{if (len .Ind) ne 0}}, Ind:{{range .Ind}} {{.}}{{end}}` +
		`{{end}}{{if (len .Msg) ne 0}}, err: '{{.Msg}}'{{end}}`
)

func init() {
	templ = template.Must(template.New("gophtu").Parse(msgFmt))
}

type arg struct {
	Expected interface{}
	Got      interface{}
	Msg      string
	Ind      []int
	ExpRes   bool
}

func assert(f func(string, ...interface{}), c bool, a arg) {
	if c != a.ExpRes {
		var b bytes.Buffer
		if err := templ.Execute(&b, a); err != nil {
			panic(err)
		}
		// based on "testing" package common.decorate func.
		_, file, line, ok := runtime.Caller(2)
		if ok {
			if idx := strings.LastIndex(file, string(os.PathSeparator)); idx >= 0 {
				file = file[idx+1:]
			}
		} else {
			file = "???"
			line = 1
		}
		// Really shouldn't do that, but wth for now.
		f("\b\b\b\b\b\b\b\b\b\b\b\b\b%s:%d: %q", file, line, b.String())
	}
}

// Wrapper for testing.Error.
func Check(t *testing.T, c bool, e, r interface{}, ind ...int) bool {
	assert(t.Errorf, c, arg{e, r, "", []int(ind), true})
	return c
}

// Wrapper for testing.Error with possibility to provide custom message.
func CheckE(t *testing.T, c bool, e, r interface{},
	msg string, ind ...int) bool {
	assert(t.Errorf, c, arg{e, r, msg, []int(ind), true})
	return c
}

// Wrapper for testing.Fatal.
func Assert(t *testing.T, c bool, e, r interface{}, ind ...int) {
	assert(t.Fatalf, c, arg{e, r, "", []int(ind), true})
}

// Wrapper for testin.Fatal with possibility to provide custom message.
func AssertE(t *testing.T, c bool, e, r interface{}, msg string, ind ...int) {
	assert(t.Fatalf, c, arg{e, r, msg, []int(ind), true})
}

// Wrapper for testing.Error with message adjusted to not equality.
func CheckFalse(t *testing.T, c bool, e interface{}, ind ...int) bool {
	assert(t.Errorf, c, arg{e, nil, "", []int(ind), false})
	return c
}

// Wrapper for testing.Error with message adjusted to not equality,
// with possibility to provide custom message.
func CheckFalseE(t *testing.T, c bool, e interface{},
	msg string, ind ...int) bool {
	assert(t.Errorf, c, arg{e, nil, msg, []int(ind), false})
	return c
}

// Wrapper for testing.Fatal with message adjusted to not equality.
func AssertFalse(t *testing.T, c bool, e interface{}, ind ...int) {
	assert(t.Fatalf, c, arg{e, nil, "", []int(ind), false})
}

// Wrapper for testing.Fatal with message adjusted to not equality,
// with possibility to provide custom message.
func AssertFalseE(t *testing.T, c bool, e interface{}, msg string, ind ...int) {
	assert(t.Fatalf, c, arg{e, nil, msg, []int(ind), false})
}
