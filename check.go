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
	// Expected represents expected value
	Expected interface{}
	// Got represents obtained value.
	Got interface{}
	// Msg represents custom error message.
	Msg string
	// Ind represents slice of indexes (useful for table tests).
	Ind []int
	// ExpRes represents expected result of comparison.
	ExpRes bool
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
		// It is stripping `file:line` prefix from original testing package.
		f("\b\b\b\b\b\b\b\b\b\b\b\b\b%s:%d: %q", file, line, b.String())
	}
}

// Check is a wrapper for testing.Error. It takes as arguments: pointer to
// *testing.T object, bool value representing result of test comparison,
// expected value, received value and optionally indexes for table tests.
func Check(t *testing.T, c bool, e, r interface{}, ind ...int) bool {
	assert(t.Errorf, c, arg{e, r, "", []int(ind), true})
	return c
}

// CheckE is a wrapper for testing.Error. It takes as arguments: pointer to
// *testing.T object, bool value representing result of test comparison,
// expected value, received value, custom error message
// and optionally indexes for table tests.
func CheckE(t *testing.T, c bool, e, r interface{},
	msg string, ind ...int) bool {
	assert(t.Errorf, c, arg{e, r, msg, []int(ind), true})
	return c
}

// Assert is a wrapper for testing.Fatal. It takes as arguments: pointer to
// *testing.T object, bool value representing result of test comparison,
// expected value, received value and optionally indexes for table tests.
func Assert(t *testing.T, c bool, e, r interface{}, ind ...int) {
	assert(t.Fatalf, c, arg{e, r, "", []int(ind), true})
}

// AssertE is a wrapper for testing.Fatal. It takes as arguments: pointer to
// *testing.T object, bool value representing result of test comparison,
// expected value, received value, custom error message
// and optionally indexes for table tests.
func AssertE(t *testing.T, c bool, e, r interface{}, msg string, ind ...int) {
	assert(t.Fatalf, c, arg{e, r, msg, []int(ind), true})
}

// CheckFalse is a wrapper for testing.Error. It is intended for
// usage when result is expected to be different than specified value.
// It takes as arguments: pointer to *testing.T object,
// bool value representing result of test comparison, expected value
// and optionally indexes for table tests.
func CheckFalse(t *testing.T, c bool, e interface{}, ind ...int) bool {
	assert(t.Errorf, c, arg{e, nil, "", []int(ind), false})
	return c
}

// CheckFalseE is a wrapper for testing.Error. It is intended for
// usage when result is expected to be different than specified value.
// It takes as arguments: pointer to *testing.T object,
// bool value representing result of test comparison, expected value,
// custom error message and optionally indexes for table tests.
func CheckFalseE(t *testing.T, c bool, e interface{},
	msg string, ind ...int) bool {
	assert(t.Errorf, c, arg{e, nil, msg, []int(ind), false})
	return c
}

// AssertFalse is a wrapper for testing.Fatal. It is intended for
// usage when result is expected to be different than specified value.
// It takes as arguments: pointer to *testing.T object,
// bool value representing result of test comparison, expected value
// and optionally indexes for table tests.
func AssertFalse(t *testing.T, c bool, e interface{}, ind ...int) {
	assert(t.Fatalf, c, arg{e, nil, "", []int(ind), false})
}

// AssertFalseE is a wrapper for testing.Fatal. It is intended for
// usage when result is expected to be different than specified value.
// It takes as arguments: pointer to *testing.T object,
// bool value representing result of test comparison, expected value,
// custom error message and optionally indexes for table tests.
func AssertFalseE(t *testing.T, c bool, e interface{}, msg string, ind ...int) {
	assert(t.Fatalf, c, arg{e, nil, msg, []int(ind), false})
}
