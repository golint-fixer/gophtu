package times

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/pblaszczyk/gophtu/asserts"
)

func Test_Timeout(t *testing.T) {
	timeouts[regexp.MustCompile("times.Test_Timeout")] = time.Minute
	timeouts[regexp.MustCompile("times.Test_Timeout2")] = time.Hour
	timeouts[regexp.MustCompile("times.Test_Timeou.*")] = 2 * time.Hour
	tT := Timeout()
	asserts.Check(t, tT == 2*time.Hour, tT, 2*time.Hour)
}

func Test_process(t *testing.T) {
	cases := []struct {
		env string
		res map[*regexp.Regexp]time.Duration
		err error
	}{
		{
			"", map[*regexp.Regexp]time.Duration{}, nil,
		},
		{
			"pack.TestThe.*=*10;l.Testu[u]+=+10ms;ur.Test_N=5h;k.Testss=+1s",
			map[*regexp.Regexp]time.Duration{
				regexp.MustCompile("pack.TestThe.*"): 10 * defaultT,
				regexp.MustCompile("l.Testu[u]+"):    defaultT + 10*time.Millisecond,
				regexp.MustCompile("ur.Test_N"):      5 * time.Hour,
				regexp.MustCompile("k.Testss"):       defaultT + time.Second,
			},
			nil,
		},
		{
			"sperpackage.TestR[R-a]=+4s",
			map[*regexp.Regexp]time.Duration{
				regexp.MustCompile("sperpackage.TestR[R-a]"): defaultT + 4*time.Second,
			},
			nil,
		},
		{
			"aljds",
			nil,
			errors.New("gophtu: invalid timeout setting: " + "aljds"),
		},
		{
			"pack.P;",
			nil,
			errors.New("gophtu: invalid timeout setting: " + "pack.P"),
		},
		{
			"[a-A]d=3s",
			nil,
			errors.New("gophtu: invalid timeout regex: " + "[a-A]d"),
		},
		{
			"sth=u",
			nil,
			errors.New("gophtu: invalid timeout setting: u"),
		},
		{
			"sd=2ms;sth=*23s",
			nil,
			errors.New("gophtu: invalid timeout setting: *23s"),
		},
	}
	for i, cas := range cases {
		timeouts = make(map[*regexp.Regexp]time.Duration)
		err := process(cas.env)
		asserts.Assert(t, (cas.err == nil) == (err == nil), cas.err, err, i)
		if cas.err != nil {
			asserts.Assert(t, strings.HasPrefix(err.Error(), cas.err.Error()),
				cas.err, err, i)
			continue
		}
		asserts.Assert(t, len(timeouts) == len(cas.res), len(cas.res),
			len(timeouts), i)
		for k, v := range cas.res {
			found := false
			for k1, v1 := range timeouts {
				if reflect.DeepEqual(*k1, *k) && reflect.DeepEqual(v, v1) {
					found = true
					break
				}
			}
			asserts.AssertE(t, found, true, found,
				fmt.Sprintf("key: %v; val: %v => map: %v", k, v, timeouts), i)
		}
	}
}
