package template

import (
	"fmt"
	"regexp"
)

func ParseTemplate(original string, vars map[string]string) (string) {
	for variable, value := range vars {
		reg := regexp.MustCompile(fmt.Sprintf(`\{\$%v\}`, variable))
		original = reg.ReplaceAllString(original, value)
	}
	return original
}