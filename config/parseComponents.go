package config

import (
	"errors"
	"regexp"
)

var trueReg = regexp.MustCompile(`(t(rue)?)|(y(es)?)|1`)
var falseReg = regexp.MustCompile(`(f(alse)?)|(n(o)?)|0`)

var errBadBool = errors.New("bad bool type")

func ParseBool(toParse string) (bool, error) {
	trueMatched := trueReg.MatchString(toParse)
	falseMatched := falseReg.MatchString(toParse)
	if !trueMatched && !falseMatched {
		return false, errBadBool
	}
	return trueMatched, nil
}