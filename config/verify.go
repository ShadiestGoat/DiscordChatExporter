package config

import "regexp"

var snowflakeVarifyReg = regexp.MustCompile(`\d{18,19}`)

func VerifySnowflake(id string) (string, bool) {
	strFound := snowflakeVarifyReg.FindString(id)
	ok := len(strFound) != 0
	return strFound, ok
}
