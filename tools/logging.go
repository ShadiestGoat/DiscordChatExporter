package tools

import (
	c "github.com/fatih/color"
)

func PanicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

var colWarn = c.New(c.Bold, c.FgHiYellow)
var colSucc = c.New(c.Bold, c.FgHiGreen)

// var colPanc = c.New(c.Bold, c.FgHiRed)

func Warn(content string) {
	colWarn.Println("Warning! " + content)
}

func Success(content string) {
	colSucc.Println(content)
}
