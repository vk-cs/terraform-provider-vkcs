package changelog

import (
	"bufio"
	"io"
	"log"
	"regexp"
	"strings"
)

var (
	versionRegexp           = regexp.MustCompile(`##? (?i:(\[UNRELEASED\]|HEAD|v?\d+\.\d+(?:\.\d+)?))(?U:.*)(\d{4}-\d{2}-\d{2})?.?$`)
	subheaderRegexp         = regexp.MustCompile(`### ([0-9A-Za-z_ ]+)`)
	changeLineRegexp        = regexp.MustCompile(`[\*|\-] (.+)`)
	changeLineRegexpWithRef = regexp.MustCompile(`[\*|\-] (.+)( \(((#[0-9]+)|(@?[[:word:]]+))\))`)
)

func matchLine(regexp *regexp.Regexp, line string) (matches []string, doesMatch bool) {
	if regexp.MatchString(line) {
		return regexp.FindStringSubmatch(line), true
	}
	return nil, false
}

func versionDateFromMatches(matches []string) string {
	var date string
	if len(matches) == 3 {
		date = matches[2]
	}
	return date
}

func parseChangelog(file io.Reader, history *Changelog) error {
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	currentHeader := ""
	currentSubHeader := ""
	var currentLine *ChangeLine
	for scanner.Scan() {
		txt := scanner.Text()
		if matches, ok := matchLine(versionRegexp, txt); ok {
			currentHeader = matches[1]
			currentSubHeader = ""
			newVersion := history.GetVersionOrCreate(currentHeader)
			newVersion.Date = versionDateFromMatches(matches)
			continue
		}

		if matches, ok := matchLine(subheaderRegexp, txt); ok {
			currentSubHeader = matches[1]
			history.GetSubsectionOrCreate(currentHeader, currentSubHeader)
			continue
		}

		if matches, ok := matchLine(changeLineRegexp, txt); ok {
			var line *ChangeLine
			if more, ok := matchLine(changeLineRegexpWithRef, txt); ok {
				// Has ref
				line = &ChangeLine{
					Summary:   more[1],
					Reference: more[3],
				}
			} else {
				// No ref
				line = &ChangeLine{
					Summary: matches[1],
				}
			}
			currentLine = line
			if currentSubHeader == "" {
				history.AddLineToVersion(currentHeader, line)
			} else {
				history.AddLineToSubsection(currentHeader, currentSubHeader, line)
			}
			continue
		} else if strings.TrimSpace(txt) != "" && currentLine != nil {
			currentLine.Summary += "\n\n" + txt
		}

	}
	if err := scanner.Err(); err != nil {
		log.Fatal("error reading history:", err)
	}
	return nil
}
