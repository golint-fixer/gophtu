package times

import (
	"errors"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	timeoutEnv = "GOPHTU_TIMEOUT"
	defaultT   = time.Millisecond
	sep        = ";"
)

var (
	timeouts = make(map[*regexp.Regexp]time.Duration)
	// Regex for timeout setting.
	timRe = regexp.MustCompile(`^([*+]?)(\d+)(ms|s|m|h)?$`)
	// Regex catching package.TestName.
	testRe = regexp.MustCompile(`^.*?([^/]+\.Test.*)$`)
)

func init() {
	if err := process(os.Getenv(timeoutEnv)); err != nil {
		panic(err.Error())
	}
}

func process(env string) error {
	if env == "" {
		return nil
	}
	for _, s := range strings.Split(env, sep) {
		t := strings.Split(s, "=")
		if len(t) != 2 {
			return errors.New("gophtu: invalid timeout setting: " + s)
		}
		r, err := regexp.Compile(t[0])
		if err != nil {
			return errors.New("gophtu: invalid timeout regex: " + t[0] +
				", err: " + err.Error())
		}
		m := timRe.FindStringSubmatch(t[1])
		if len(m) != 4 {
			return errors.New("gophtu: invalid timeout setting: " + t[1])
		}
		u, err := strconv.ParseUint(m[2], 10, 64)
		if err != nil {
			return errors.New("gophtu: " + err.Error())
		}
		if (m[1] == "*" && m[3] != "") || (m[1] != "*" && m[3] == "") {
			return errors.New("gophtu: invalid timeout setting: " + t[1])
		}
		c := map[string]time.Duration{
			"ms": time.Millisecond,
			"s":  time.Second,
			"m":  time.Minute,
			"h":  time.Hour,
		}
		switch m[1] {
		case "*":
			timeouts[r] = defaultT * time.Duration(u)
		case "+":
			timeouts[r] = defaultT + time.Duration(u)*c[m[3]]
		default:
			timeouts[r] = time.Duration(u) * c[m[3]]
		}
	}
	return nil
}

// Timeout returns default timeout, if one is not explicitly configured
// for test, or max matching timeout otherwise.
// Timeout uses GOPHTU_TIMEOUT environment variable to check for
// preconfigured timeouts.
// Syntax for var is: {op}{testregex}{tu};{op}{testregex}{tu}
// where {op} is one of "+", "*", "", {testregex} is regex for test
// in format package.TestName, {tu} is time unit from: "ms", "s", "m", "h", "".
// {tu} must be set for {op} equal "*" and can't be set otherwise.
func Timeout() time.Duration {
	pc := make([]uintptr, 10)
	var str string
	runtime.Callers(0, pc)
	for _, pc := range pc {
		if f := runtime.FuncForPC(pc); f != nil {
			m := testRe.FindStringSubmatch(f.Name())
			if len(m) > 1 {
				str = m[1]
				break
			}
		}
	}
	max := defaultT
	for r, t := range timeouts {
		if r.MatchString(str) {
			if t > max {
				max = t
			}
		}
	}
	return max
}
