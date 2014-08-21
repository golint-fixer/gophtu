package gophtu

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func mock(sl *string) func(s string, i ...interface{}) {
	return func(s string, i ...interface{}) {
		*sl = fmt.Sprintf(s, i...)
	}
}

func runAssert(f func(func(string, ...interface{}), bool, arg),
	f2 func(*string) func(string, ...interface{}), s *string, b bool, a arg) int {
	f(f2(s), b, a)
	return lineNo()
}

func lineNo() int {
	_, _, line, ok := runtime.Caller(2)
	if !ok {
		line = 1
	}
	return line
}

func Test_assert(t *testing.T) {
	cfg := []struct {
		err string
		c   bool
		a   arg
	}{
		{"", true, arg{ExpRes: true}},
		{"check_test.go:%d: \"want: '2', got: '3'\"", false,
			arg{ExpRes: true, Expected: 2, Got: 3}},
		{"check_test.go:%d: \"want: '3', got: '4', err: 'ziemniak'\"", false,
			arg{ExpRes: true, Expected: 3, Got: 4, Msg: "ziemniak"}},
		{"check_test.go:%d: \"want: '3', got: '4', Ind: 2 4, err: 'ziemniak'\"",
			false, arg{ExpRes: true, Expected: 3, Got: 4,
				Msg: "ziemniak", Ind: []int{2, 4}}},
		{"check_test.go:%d: \"want: '3', got: '4', Ind: 3 5\"", false,
			arg{ExpRes: true, Expected: 3, Got: 4, Ind: []int{3, 5}}},
		{"check_test.go:%d: \"want: result different than '5', Ind: 2 3\"", true,
			arg{ExpRes: false, Expected: 5, Ind: []int{2, 3}}},
		{"check_test.go:%d: \"want: result different than '5'\"", true,
			arg{ExpRes: false, Expected: 5}},
		{"check_test.go:%d: \"want: result different than '5', err: " +
			"'I am the error'\"", true, arg{ExpRes: false, Expected: 5,
			Msg: "I am the error"}},
		{"check_test.go:%d: \"want: result different than '5', Ind: 2 3 4, " +
			"err: 'mSg a'\"", true, arg{ExpRes: false, Expected: 5,
			Msg: "mSg a", Ind: []int{2, 3, 4}}},
	}

	for i, cfg := range cfg {
		var s string
		n := runAssert(assert, mock, &s, cfg.c, cfg.a)
		st := strings.TrimLeft(s, "\b")
		err := cfg.err
		if err != "" {
			err = fmt.Sprintf(err, n)
		}
		if !strings.HasPrefix(st, err) {
			t.Errorf("expected '%v' starts with '%v' (%d)", st, err, i)
		}
		if st != err {
			tokens := strings.Split(st, ":")
			if len(tokens) == 0 {
				t.Fatalf("expecte len(tokens)=0; got %d", len(tokens))
			}
			if !filepath.IsAbs(filepath.FromSlash(tokens[0])) {
				t.Errorf("expected filepath.Abs(tokens[0])=true")
			}
		}
	}
}
