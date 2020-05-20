package html2text

import (
	"bytes"
	"regexp"
	"strconv"
	"strings"
)

const (
	WIN_LBR  = ""
	UNIX_LBR = ""
)

var lbr = WIN_LBR
var badTagnamesRE = regexp.MustCompile(`^(head|script|style)($|\s*)`)
var linkTagRE = regexp.MustCompile(`a.*href=('([^']*?)'|"([^"]*?)")`)
var badLinkHrefRE = regexp.MustCompile(`#|javascript:`)
var headersRE = regexp.MustCompile(`^(\/)?h[1-6]`)
var numericEntityRE = regexp.MustCompile(`^#([0-9]+)$`)

func parseHTMLEntity(entName string) (string, bool) {
	if r, ok := entity[entName]; ok {
		return string(r), true
	}

	if match := numericEntityRE.FindStringSubmatch(entName); len(match) == 2 {
		digits := match[1]
		n, err := strconv.Atoi(digits)
		if err == nil && (n == 9 || n == 10 || n == 13 || n > 31) {
			return string(rune(n)), true
		}
	}

	return "", false
}

// SetUnixLbr with argument true sets Unix-style line-breaks in output ("\n")
// with argument false sets Windows-style line-breaks in output ("\r\n", the default)
func SetUnixLbr(b bool) {
	if b {
		lbr = UNIX_LBR
	} else {
		lbr = WIN_LBR
	}
}

// HTMLEntitiesToText decodes HTML entities inside a provided
// string and returns decoded text
func HTMLEntitiesToText(htmlEntsText string) string {
	outBuf := bytes.NewBufferString("")
	inEnt := false

	for i, r := range htmlEntsText {
		switch {
		case r == ';' && inEnt:
			inEnt = false
			continue

		case r == '&': //possible html entity
			entName := ""
			isEnt := false

			// parse the entity name - max 10 chars
			chars := 0
			for _, er := range htmlEntsText[i+1:] {
				if er == ';' {
					isEnt = true
					break
				} else {
					entName += string(er)
				}

				chars++
				if chars == 10 {
					break
				}
			}

			if isEnt {
				if ent, isEnt := parseHTMLEntity(entName); isEnt {
					outBuf.WriteString(ent)
					inEnt = true
					continue
				}
			}
		}

		if !inEnt {
			outBuf.WriteRune(r)
		}
	}

	return outBuf.String()
}

func writeSpace(outBuf *bytes.Buffer) {
	bts := outBuf.Bytes()
	if len(bts) > 0 && bts[len(bts)-1] != ' ' {
		outBuf.WriteString(" ")
	}
}

func sReplacer(text string) string {
	rsbt := regexp.MustCompile(`>\s+<`)
	t := rsbt.ReplaceAllString(text, "><")
	re := regexp.MustCompile(`\x{000D}\x{000A}|[\x{000A}\x{000B}\x{000C}\x{000D}\x{0085}\x{2028}\x{2029}]`)
	t = re.ReplaceAllString(t, "")
	var r = strings.NewReplacer("\t", "", "\f", "")
	t = r.Replace(t)
	return t
}

// HTML2Text converts html into a text form
func HTML2Text(html string) string {
	html = sReplacer(html)
	tagStart := 0
	inEnt := false
	badTagStackDepth := 0 // if == 1 it means we are inside <head>...</head>
	shouldOutput := true

	outBuf := bytes.NewBufferString("")

	for i, r := range html {

		switch {
		// skip new lines and spaces adding a single space if not there yet
		case r <= 0xD, r == 0x85, r == 0x2028, r == 0x2029, // new lines
			r == ' ', r >= 0x2008 && r <= 0x200B: // spaces
			writeSpace(outBuf)
			continue

		case r == ';' && inEnt: // end of html entity
			inEnt = false
			shouldOutput = true
			continue

		case r == '&' && shouldOutput: // possible html entity
			entName := ""
			isEnt := false

			// parse the entity name - max 10 chars
			chars := 0
			for _, er := range html[i+1:] {
				if er == ';' {
					isEnt = true
					break
				} else {
					entName += string(er)
				}

				chars++
				if chars == 10 {
					break
				}
			}

			if isEnt {
				if ent, isEnt := parseHTMLEntity(entName); isEnt {
					outBuf.WriteString(ent)
					inEnt = true
					shouldOutput = false
					continue
				}
			}

		case r == '<': // start of a tag
			tagStart = i + 1
			shouldOutput = false
			continue

		case r == '>': // end of a tag
			shouldOutput = true
			tag := html[tagStart:i]
			tagNameLowercase := strings.ToLower(tag)

			if tagNameLowercase == "/ul" {
				outBuf.WriteString(lbr)
			} else if tagNameLowercase == "li" || tagNameLowercase == "li/" {
				outBuf.WriteString(" ")
			} else if headersRE.MatchString(tagNameLowercase) {
				outBuf.WriteString(" ")
			} else if tagNameLowercase == "br" || tagNameLowercase == "br/" {
				// new line
				outBuf.WriteString(lbr)
			} else if tagNameLowercase == "p" {
				outBuf.WriteString(lbr + lbr)
			} else if tagNameLowercase == "/p" {
				outBuf.WriteString(" ")
			} else if badTagnamesRE.MatchString(tagNameLowercase) {
				// unwanted block
				badTagStackDepth++
			} else if len(tagNameLowercase) > 0 && tagNameLowercase[0] == '/' &&
				badTagnamesRE.MatchString(tagNameLowercase[1:]) {
				// end of unwanted block
				badTagStackDepth--
			}
			continue

		} // switch end

		if shouldOutput && badTagStackDepth == 0 && !inEnt {
			outBuf.WriteRune(r)
		}
	}

	tOut := outBuf.String()
	tOut = strings.TrimSpace(tOut)
	return tOut
}
