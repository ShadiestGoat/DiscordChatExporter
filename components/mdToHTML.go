package components

import (
	"fmt"
	"regexp"
)

var rReg = regexp.MustCompile(`\r`)

var mdSyntax = map[string]string{
	"\\*": "*",
	"_":   "_",
	"\\|": "|",
	"`":   "`",
	"~":   "~",
}

var mdSpreadInline = createMdInfo("`", "`", 1, NB_Spread, "", "")

var (
	mdBold       = createMdInfo(`\*`, "*", 2, NB_Spread, "b", "")
	mdUnderline  = createMdInfo(`_`, "_", 2, NB_Spread, "u", "")
	mdItal1      = createMdInfo(`\*`, "*", 1, NB_Ignore, "i", "")
	mdItal2      = createMdInfo(`_`, "_", 1, NB_Ignore, "i", "")
	mdStrike     = createMdInfo(`~`, "~", 2, NB_Stop, "s", "")
	mdSpoiler    = createMdInfo(`\|`, "|", 2, NB_Stop, "span", "spoiler")
	mdInlineCode = createMdInfo("`", "`", 1, NB_Ignore, "span", "inline-code")
)

var mdMap = []mdInfo{
	mdBold,
	mdUnderline,
	mdItal1,
	mdItal2,
	mdStrike,
	mdSpoiler,
}

type mdInfo struct {
	Reg *regexp.Regexp

	Char byte
	Size int

	NextBehavior NextBehavior

	HtmlOpen  string
	HtmlClose string
}

// Caller is responsibe to make sure there is a match
func (m mdInfo) ReplaceAll(content *string) {
	skipI := 0

	for skipI != len(m.Reg.FindAllStringIndex(*content, skipI+1)) {
		locs := m.Reg.FindAllStringIndex(*content, skipI+1)
		loc := locs[skipI]

		switch m.NextBehavior {
		case NB_Ignore:
			if loc[1] != len(*content)-1 {
				loc[1] -= 1
			}
		}

		locInner := []int{loc[0] + m.Size, loc[1] - m.Size}

		*content = (*content)[:loc[0]] + m.HtmlOpen + (*content)[locInner[0]:locInner[1]] + m.HtmlClose + (*content)[loc[1]:]
	}

}

// func (m mdInfo)

type NextBehavior int8

const (
	// The things must put on the out most characters of the group
	NB_Spread NextBehavior = iota
	// Stop at the first match, ie.
	//
	// ~~~asdasdad~~~
	//
	// ^^_________^^_
	NB_Stop
	// Ignores the match if there is another one of it's symbols after it. ie.
	//
	// _aa_: ok
	//
	// __aa_: ok
	//
	// _^__^
	//
	// _aa__: not ok
	NB_Ignore
)

func MDToHTML(content string) (string, bool) {
	inlineCodeMap := map[int]string{}
	escapedSyntaxMap := map[int]string{}
	icI := 0
	esI := 0

	for syntax := range mdSyntax {
		reg := regexp.MustCompile(`\\` + syntax)
		for reg.MatchString(content) {
			loc := reg.FindStringIndex(content)
			escapedSyntaxMap[esI] = syntax
			content = content[:loc[0]] + fmt.Sprintf("--\tes%v--", esI) + content[loc[1]:]
			esI++
		}
	}

	// inline code index
	for mdSpreadInline.Reg.MatchString(content) {
		curLoc := mdSpreadInline.Reg.FindStringIndex(content)
		inlineCodeMap[icI] = content[curLoc[0]:curLoc[1]]
		content = fmt.Sprintf("%v--\tic%v--%v", content[:curLoc[0]], icI, content[curLoc[1]:])
		icI++
	}

	bq := false
	if len(content) > 5 {
		if content[:5] == "&gt; " {
			// Block quote
			content = content[5:]
			bq = true
		}
	}

	for _, mdRule := range mdMap {
		mdRule.ReplaceAll(&content)
	}

	for i, str := range inlineCodeMap {
		reg := regexp.MustCompile(fmt.Sprintf(`--\tic%v--`, i))
		content = reg.ReplaceAllString(content, str)
	}

	for i, str := range escapedSyntaxMap {
		reg := regexp.MustCompile(fmt.Sprintf(`--\tes%v--`, i))
		content = reg.ReplaceAllString(content, mdSyntax[str])
	}

	return content, bq
}

func createMdInfo(regEscapeChar string, oneChar string, size int, nextBehavior NextBehavior, htmlTag string, cssClass string) mdInfo {
	htmlAttr := ""

	if cssClass != "" {
		htmlAttr = " class=\"" + cssClass + "\""
	}

	return mdInfo{
		Reg:          createEnclosed(regEscapeChar, size, nextBehavior),
		HtmlOpen:     "<" + htmlTag + htmlAttr + ">",
		HtmlClose:    fmt.Sprintf(`</%v>`, htmlTag),
		Char:         []byte(oneChar)[0],
		NextBehavior: nextBehavior,
		Size:         size,
	}
}

func createEnclosed(charReg string, size int, nextBehavior NextBehavior) *regexp.Regexp {
	strReg := ""
	switch nextBehavior {
	case NB_Ignore:
		totalReg := ""
		for i := 0; i < size; i++ {
			totalReg += charReg
		}
		strReg = fmt.Sprintf(`%v[^%v]+?%v([^%v]|$)`, totalReg, charReg, totalReg, charReg)
	case NB_Spread:
		strReg = fmt.Sprintf(`%v{%v,}.+?%v{%v,}`, charReg, size, charReg, size)
	case NB_Stop:
		strReg = fmt.Sprintf(`%v{%v,}.+?%v{%v}`, charReg, size, charReg, size)
	}

	return regexp.MustCompile(strReg)
}
