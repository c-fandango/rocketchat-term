package main

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/c-fandango/rocketchat-term/utils"
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
	const roomWidth = 24
	const userWidth = 14
	const indentWidth = 7
	const newLineMarkerWidth = 14
	const roomNameMaxWidth = 23
	const contentIndent = timeWidth + roomWidth + userWidth + indentWidth + 2

	resetColour := "\033[0m"
	blackText := "\033[38;5;0m"
	ticketColour := "\033[38;5;39m" // deepSkyBlue1

	replacePatterns := map[string]string{
		`( |^)(@[^\s]+)`:    fmt.Sprintf("${1}%s ${2} %s", config.notifyColour, resetColour),
		`( |^)(#\d{6})`:     fmt.Sprintf("${1}%s${2}%s", ticketColour, resetColour),
		"```((.|\\n)+?)```": fmt.Sprintf("%s${1}%s", config.codeColour, resetColour),
		`(\n)`:              "\n" + strings.Repeat(" ", contentIndent),
		`(\z)`:              resetColour,
	}

	// weird but has to be executed after the other code highlighting regex (negative lookarounds are not supported)
	replaceCodeline := map[string]string{
		"`((.|\\n)+?)`": config.codeColour + "${1}" + resetColour,
	}

	userColour := config.textColours[len(user)%(len(config.textColours)-1)]
	roomColour := config.bgColours[len(room)%(len(config.bgColours)-1)] + blackText

	user = makeShortName(user)

	ts := time.UnixMilli(int64(timestamp))
	timePretty := ts.Format(time.Kitchen)
	timePretty = utils.PadRight(timePretty, " ", timeWidth)

	roomNameMaxIndex := utils.MinInt(roomNameMaxWidth, len(room))
	room = room[:roomNameMaxIndex]
	roomFmt := roomColour + " " + room + " " + resetColour
	userFmt := userColour + user + resetColour

	roomFmtWidth := roomWidth + len(roomFmt) - len(room)
	userFmtWidth := userWidth + len(userFmt) - len(user)

	newLine := strings.Repeat(" ", indentWidth) + timePretty + utils.PadRight(roomFmt, " ", roomFmtWidth) + utils.PadRight(userFmt, " ", userFmtWidth) + fmtContent(fmtContent(content, replacePatterns), replaceCodeline)

	newLine = strings.Repeat("-", newLineMarkerWidth) + "\n" + newLine

	fmt.Println(newLine)
}
