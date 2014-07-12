package gophtu

import (
	"fmt"
	"strings"
	"testing"
)

func mock(sl *string) func(s string, i ...interface{}) {
	return func(s string, i ...interface{}) {
		*sl = fmt.Sprintf(s, i...)
	}
}

func runAssert(f func(func(string, ...interface{}), bool, arg),
	f2 func(*string) func(string, ...interface{}), s *string, b bool, a arg) {
	f(f2(s), b, a)
}

func Test_assert(t *testing.T) {
	// TODO: drop format dependency
	cfg := []struct {
		err string
		c   bool
		a   arg
	}{
		{"", true, arg{ExpRes: true}},
		{"check_test.go:40: \"want: '2', got: '3'\"", false,
			arg{ExpRes: true, Expected: 2, Got: 3}},
		{"check_test.go:40: \"want: '3', got: '4', err: 'ziemniak'\"", false,
			arg{ExpRes: true, Expected: 3, Got: 4, Msg: "ziemniak"}},
		{"check_test.go:40: \"want: '3', got: '4', Ind: 2 4, err: 'ziemniak'\"",
			false, arg{ExpRes: true, Expected: 3, Got: 4,
				Msg: "ziemniak", Ind: []int{2, 4}}},
		{"check_test.go:40: \"want: '3', got: '4', Ind: 3 5\"", false,
			arg{ExpRes: true, Expected: 3, Got: 4, Ind: []int{3, 5}}},
	}

	for i, cfg := range cfg {
		var s string
		runAssert(assert, mock, &s, cfg.c, cfg.a)
		st := strings.TrimLeft(s, "\b")
		if cfg.err != st {
			func() {
				t.Errorf("expected '%v' equal '%v' (%d)", cfg.err, st, i)
			}()
		}
	}
}
