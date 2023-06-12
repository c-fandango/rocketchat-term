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

	var contentIndent = config.timeWidth + config.roomWidth + config.userWidth + config.indentWidth + 2
	resetColour := "\033[0m"

	replacePatterns := map[string]string{
		`( |^)(@[^\s]+)`:    fmt.Sprintf("${1}%s ${2} %s", config.notifyColour, resetColour),
		`( |^)(#\d{6})`:     fmt.Sprintf("${1}%s${2}%s", config.ticketColour, resetColour),
		"```((.|\\n)+?)```": fmt.Sprintf("%s${1}%s", config.codeColour, resetColour),
		`(\n)`:              "\n" + strings.Repeat(" ", contentIndent),
		`(\z)`:              resetColour,
	}

	// weird but has to be executed after the other code highlighting regex (negative lookarounds are not supported)
	replaceCodeline := map[string]string{
		"`((.|\\n)+?)`": config.codeColour + "${1}" + resetColour,
	}

	userTextColour := config.userTextColours[len(user)%len(config.userTextColours)]
	roomTextColour := config.roomTextColours[len(room)%len(config.roomTextColours)]
	userBgColour := config.userBgColours[len(user)%len(config.userBgColours)]
	roomBgColour := config.roomBgColours[len(room)%len(config.roomBgColours)]

	userColour := userBgColour + userTextColour
	roomColour := roomBgColour + roomTextColour

	user = makeShortName(user)

	ts := time.UnixMilli(int64(timestamp))
	timePretty := ts.Format(time.Kitchen)
	timePretty = utils.PadRight(timePretty, " ", config.timeWidth)

	roomNameMaxIndex := utils.MinInt(config.roomNameMaxWidth, len(room))
	room = room[:roomNameMaxIndex]
	roomFmt := roomColour + " " + room + " " + resetColour
	userFmt := userColour + user + resetColour

	roomFmtWidth := config.roomWidth + len(roomFmt) - len(room)
	userFmtWidth := config.userWidth + len(userFmt) - len(user)

	newLine := strings.Repeat(" ", config.indentWidth) + timePretty + utils.PadRight(roomFmt, " ", roomFmtWidth) + utils.PadRight(userFmt, " ", userFmtWidth) + fmtContent(fmtContent(content, replacePatterns), replaceCodeline)

	newLine = strings.Repeat("-", config.newLineMarkerWidth) + "\n" + newLine

	fmt.Println(newLine)
}
