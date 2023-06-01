package display

import (
	"fmt"
	"github.com/c-fandango/rocketchat-term/utils"
	"time"
)

func PrintMessage(room string, user string, content string, timestamp int) {

	colourTextCodes := []string{
		"31", "32", "33", "34", "35", "36",
	}
	colourBackgroundCodes := []string{"41", "42", "43", "44", "45", "46"}

	resetColour := "\033[0m"
	blackText := "\033[30m"

	userColour := "\033[" + colourTextCodes[len(user)%(len(colourTextCodes)-1)] + "m"
	roomColour := "\033[" + colourBackgroundCodes[len(room)%(len(colourBackgroundCodes)-1)] + "m"

	user = userColour + user + resetColour
	room = roomColour + blackText + " " + room + " " + resetColour

	ts := time.UnixMilli(int64(timestamp))
	timePretty := ts.Format(time.Kitchen)

	timePretty = utils.PadLeft(timePretty, " ", 14)

	fmt.Println(timePretty, "     ", room, "      ", user, "     ", content)
}
