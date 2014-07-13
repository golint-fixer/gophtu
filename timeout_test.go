package gophtu

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"testing"
	"time"
)

func Test_Timeout(t *testing.T) {
	timeouts[regexp.MustCompile("gophtu.Test_Timeout")] = time.Minute
	timeouts[regexp.MustCompile("gophtu.Test_Timeout2")] = time.Hour
	timeouts[regexp.MustCompile("gophtu.Test_Timeou.*")] = 2 * time.Hour
	tT := Timeout()
	Check(t, tT == 2*time.Hour, tT, 2*time.Hour)
}

func Test_process(t *testing.T) {
	cfg := []struct {
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
	for i := range cfg {
		timeouts = make(map[*regexp.Regexp]time.Duration)
		err := process(cfg[i].env)
		Assert(t, (cfg[i].err == nil) == (err == nil), cfg[i].err, err, i)
		if cfg[i].err != nil {
			Assert(t, strings.HasPrefix(err.Error(), cfg[i].err.Error()),
				cfg[i].err, err, i)
			continue
		}
		Assert(t, len(timeouts) == len(cfg[i].res), len(cfg[i].res),
			len(timeouts), i)
		for k, v := range cfg[i].res {
			found := false
			for k1, v1 := range timeouts {
				if reflect.DeepEqual(*k1, *k) && reflect.DeepEqual(v, v1) {
					found = true
					break
				}
			}
			AssertE(t, found, true, found,
				fmt.Sprintf("key: %v; val: %v => map: %v", k, v, timeouts), i)
		}
	}
}
