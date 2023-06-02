package main

import (
	"fmt"
	"github.com/c-fandango/rocketchat-term/utils"
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

func makeRoomName(input []string) string {
	if len(input) == 1 {
		return input[0]
	}

	initials := make([]string, len(input))

	for i, name := range input {
		initials[i] = makeInitials(name, ".")
	}

	return strings.Join(initials, ", ")

}

func fmtContent(content string, indent int) string {

	lines := strings.Split(content, "\n")
	for i, line := range lines {
		if i != 0 {
			line = strings.Repeat(" ", indent) + line
		}

		lines[i] = strings.ReplaceAll(line, "```", "")
	}

	return strings.Join(lines, "\n")

}

func printMessage(room string, user string, content string, timestamp int) {

	const timeWidth = 15
	const roomWidth = 24
	const userWidth = 14
	const indentWidth = 7
	const contentIndent = timeWidth + roomWidth + userWidth + indentWidth + 2

	colourCodes256 := []string{
		"2", "5", "6", "20", "40", "87", "93", "130", "158", "171", "214", "220",
	}

	resetColour := "\033[0m"
	blackText := "\033[38;5;0m"

	userColour := "\033[38;5;" + colourCodes256[len(user)%(len(colourCodes256)-1)] + "m"
	roomColour := "\033[48;5;" + colourCodes256[len(room)%(len(colourCodes256)-1)] + "m" + blackText

	user = makeShortName(user)

	ts := time.UnixMilli(int64(timestamp))
	timePretty := ts.Format(time.Kitchen)
	timePretty = utils.PadRight(timePretty, " ", timeWidth)

	roomFmt := roomColour + " " + room + " " + resetColour
	userFmt := userColour + user + resetColour

	roomFmtWidth := roomWidth + len(roomFmt) - len(room)
	userFmtWidth := userWidth + len(userFmt) - len(user)

	newLine := strings.Repeat(" ", indentWidth) + timePretty + utils.PadRight(roomFmt, " ", roomFmtWidth) + utils.PadRight(userFmt, " ", userFmtWidth) + fmtContent(content, contentIndent)

	fmt.Println(newLine)
}
