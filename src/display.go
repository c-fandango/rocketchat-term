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

func fmtContent(content string, indent int, codeColour string) string {
	resetColour := "\033[0m"

        content = utils.ReplaceEveryOther(content, "`", resetColour+codeColour)
        content = strings.ReplaceAll(content, "`", resetColour)

	lines := strings.Split(content, "\n")
	for i:=1; i<len(lines); i++ {
			lines[i] = strings.Repeat(" ", indent) + lines[i]
	}

	return strings.Join(lines, "\n")
}

func printMessage(room string, user string, content string, timestamp int) {

	const timeWidth = 15
	const roomWidth = 24
	const userWidth = 14
	const indentWidth = 7
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
		"\033[38;5;169m", // hotPink2
		"\033[48;5;171m", // mediumOrchid1
		"\033[48;5;214m", // orange1
		"\033[48;5;220m", // gold1
	}

	resetColour := "\033[0m"
	blackText := "\033[38;5;0m"
	codeColour := "\033[38;5;186m"

	userColour := fgColourCodes256[len(user)%(len(fgColourCodes256)-1)]
	roomColour := bgColourCodes256[len(room)%(len(bgColourCodes256)-1)] + blackText

	user = makeShortName(user)

	ts := time.UnixMilli(int64(timestamp))
	timePretty := ts.Format(time.Kitchen)
	timePretty = utils.PadRight(timePretty, " ", timeWidth)

	roomFmt := roomColour + " " + room + " " + resetColour
	userFmt := userColour + user + resetColour

	roomFmtWidth := roomWidth + len(roomFmt) - len(room)
	userFmtWidth := userWidth + len(userFmt) - len(user)

	newLine := resetColour + strings.Repeat(" ", indentWidth) + timePretty + utils.PadRight(roomFmt, " ", roomFmtWidth) + utils.PadRight(userFmt, " ", userFmtWidth) + fmtContent(content, contentIndent, codeColour)

	fmt.Println(newLine)
}
