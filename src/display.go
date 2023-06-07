package main

import (
	"fmt"
	"github.com/c-fandango/rocketchat-term/utils"
	"regexp"
	"strings"
	"time"
)

func makeInitials(name string, delimiter string) string {

	var initials string
	names := strings.Split(name, delimiter)

	for _, name := range names {
		initials += string(name[0])
	}

	return strings.ToUpper(initials)
}

func makeShortName(name string) string {

	names := strings.Split(name, " ")

	if len(names) == 1 {
		return name
	}

	shortName := names[0] + " " + makeInitials(names[1], "-")

	return shortName
}

func initialiseNames(input []string) string {
	if len(input) == 1 {
		return input[0]
	}

	initials := make([]string, len(input))

	for i, name := range input {
		initials[i] = makeInitials(name, ".")
	}

	return strings.Join(initials, ", ")
}

func fmtContent(content string, replacePatterns map[string]string) string {

	for pattern, replace := range replacePatterns {
		reg := regexp.MustCompile(pattern)
		content = reg.ReplaceAllString(content, replace)
	}

	return content
}

func printMessage(room string, user string, content string, timestamp int) {

	const timeWidth = 15
	const roomWidth = 23
	const userWidth = 14
	const indentWidth = 7
	const newLineMarkerWidth = 14
	const roomNameMaxWidth = 24
	const contentIndent = timeWidth + roomWidth + userWidth + indentWidth + 2

	fgColourCodes256 := []string{
		"\033[38;5;2m",   // green
		"\033[38;5;5m",   // purple
		"\033[38;5;6m",   // teal
		"\033[38;5;40m",  // green3
		"\033[38;5;87m",  // darkSlateGray2
		"\033[38;5;130m", // darkOrange3
		"\033[38;5;148m", // yellow3
		"\033[38;5;158m", // darkSeaGreen1
		"\033[38;5;169m", // hotPink2
		"\033[38;5;171m", // mediumOrchid1
		"\033[38;5;214m", // orange1
		"\033[38;5;220m", // gold1
	}

	bgColourCodes256 := []string{
		"\033[48;5;2m",   // green
		"\033[48;5;5m",   // purple
		"\033[48;5;6m",   // teal
		"\033[48;5;40m",  // green3
		"\033[48;5;87m",  // darkSlateGray2
		"\033[48;5;130m", // darkOrange3
		"\033[48;5;148m", // yellow3
		"\033[48;5;158m", // darkSeaGreen1
		"\033[48;5;169m", // hotPink2
		"\033[48;5;171m", // mediumOrchid1
		"\033[48;5;214m", // orange1
		"\033[48;5;220m", // gold1
	}

	resetColour := "\033[0m"
	blackText := "\033[38;5;0m"
	codeColour := "\033[38;5;186m"   // lightGoldenrod2
	notifyColour := "\033[48;5;160m" // red3
	ticketColour := "\033[38;5;39m"  // deepSkyBlue1

	replacePatterns := map[string]string{
		`( |^)(@[^\s]+)`:    "${1}" + notifyColour + " ${2} " + resetColour,
		`( |^)(#\d{6})`:     "${1}" + ticketColour + "${2}" + resetColour,
		"```((.|\\n)+?)```": codeColour + "${1}" + resetColour,
		"`((.|\\n)+?)`":     codeColour + "${1}" + resetColour,
		`(\n)`:              "\n" + strings.Repeat(" ", contentIndent),
		`(\z)`:              resetColour,
	}

	userColour := fgColourCodes256[len(user)%(len(fgColourCodes256)-1)]
	roomColour := bgColourCodes256[len(room)%(len(bgColourCodes256)-1)] + blackText

	user = makeShortName(user)

	ts := time.UnixMilli(int64(timestamp))
	timePretty := ts.Format(time.Kitchen)
	timePretty = utils.PadRight(timePretty, " ", timeWidth)

	roomNameMaxIndex := utils.MinInt(roomNameMaxWidth-1, len(room)-1)
	room = room[:roomNameMaxIndex]
	roomFmt := roomColour + " " + room + " " + resetColour
	userFmt := userColour + user + resetColour

	roomFmtWidth := roomWidth + len(roomFmt) - len(room)
	userFmtWidth := userWidth + len(userFmt) - len(user)

	newLine := strings.Repeat(" ", indentWidth) + timePretty + utils.PadRight(roomFmt, " ", roomFmtWidth) + utils.PadRight(userFmt, " ", userFmtWidth) + fmtContent(content, replacePatterns)

	newLine = strings.Repeat("-", newLineMarkerWidth) + "\n" + newLine

	fmt.Println(newLine)
}
