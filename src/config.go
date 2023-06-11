package main

import (
	"fmt"
	"os"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"

	"github.com/c-fandango/rocketchat-term/utils"
)

var defaultCols = []string{
	"2",   // green
	"5",   // purple
	"6",   // teal
	"40",  // green3
	"87",  // darkSlateGray2
	"130", // darkOrange3
	"148", // yellow3
	"158", // darkSeaGreen1
	"169", // hotPink2
	"171", // mediumOrchid1
	"214", // orange1
	"220", // gold1
}

var black = []string{"\033[38;5;0m"}
var nothing = []string{"\033[0m"}

const defaultCode = "\033[38;5;186m"
const defaultNotify = "\033[48;5;160m"

type configSchema struct {
	userTextColours    []string
	userBgColours      []string
	roomTextColours    []string
	roomBgColours      []string
	codeColour         string
	notifyColour       string
	timeWidth          int
	roomWidth          int
	userWidth          int
	indentWidth        int
	newLineMarkerWidth int
	roomNameMaxWidth   int
	debug              bool
}

func hexToAnsi(prefix string) func(string) string {
	return func(code string) string {
		r, g, b, err := utils.HexToRGB(code)
		if err != nil {
			panic("invalid hex code")
		}
		return fmt.Sprintf("%s;%d;%d;%dm", prefix, r, g, b)
	}
}

func numToAnsi(prefix string) func(string) string {
	return func(code string) string {
		return fmt.Sprintf("%s;%sm", prefix, code)
	}
}

func (c *configSchema) loadConf(path string) {
	var k = koanf.New(".")

	if _, err := os.Stat(path); err == nil {
		k.Load(file.Provider(path), yaml.Parser())
	}

	// read logging opts
	c.debug = k.Bool("logging.debug")

	// set indent defaults
	c.timeWidth = 15
	c.roomWidth = 24
	c.userWidth = 14
	c.indentWidth = 7
	c.newLineMarkerWidth = 14
	c.roomNameMaxWidth = 23

	// read indent opts
	if n := k.Int("spacing.time"); n != 0 {
		c.timeWidth = n
	}
	if n := k.Int("spacing.room"); n != 0 {
		c.roomWidth = n
	}
	if n := k.Int("spacing.user"); n != 0 {
		c.userWidth = n
	}
	if n := k.Int("spacing.indent"); n != 0 {
		c.indentWidth = n
	}
	if n := k.Int("spacing.marker"); n != 0 {
		c.newLineMarkerWidth = n
	}
	if n := k.Int("spacing.room_max_length"); n != 0 {
		c.roomNameMaxWidth = n
	}

	// set colour defaults
	c.userTextColours = utils.MapperStr(defaultCols, numToAnsi("\033[38;5"))
	c.userBgColours = nothing
	c.roomTextColours = black
	c.roomBgColours = utils.MapperStr(defaultCols, numToAnsi("\033[48;5"))
	c.codeColour = defaultCode
	c.notifyColour = defaultNotify

	// read colour opts
	if len(k.Strings("colours.user_text")) != 0 {
		c.userTextColours = utils.MapperStr(k.Strings("colours.user_text"), hexToAnsi("\033[38;2"))
	} else if len(k.Strings("colours256.user_text")) != 0 {
		c.userTextColours = utils.MapperStr(k.Strings("colours256.user_text"), numToAnsi("\033[38;5"))
	}

	if len(k.Strings("colours.user_highlight")) != 0 {
		c.userBgColours = utils.MapperStr(k.Strings("colours.highlight"), hexToAnsi("\033[48;2"))
	} else if len(k.Strings("colours256.user_highlight")) != 0 {
		c.userBgColours = utils.MapperStr(k.Strings("colours256.highlight"), numToAnsi("\033[48;5"))
	}

	if len(k.Strings("colours.room_text")) != 0 {
		c.roomTextColours = utils.MapperStr(k.Strings("colours.room_text"), hexToAnsi("\033[38;2"))
	} else if len(k.Strings("colours256.room_text")) != 0 {
		c.roomTextColours = utils.MapperStr(k.Strings("colours256.room_text"), numToAnsi("\033[38;5"))
	}

	if len(k.Strings("colours.room_highlight")) != 0 {
		c.roomBgColours = utils.MapperStr(k.Strings("colours.highlight"), hexToAnsi("\033[48;2"))
	} else if len(k.Strings("colours256.room_highlight")) != 0 {
		c.roomBgColours = utils.MapperStr(k.Strings("colours256.highlight"), numToAnsi("\033[48;5"))
	}

	if len(k.String("colours.code")) != 0 {
		c.codeColour = hexToAnsi("\033[38;2")(k.String("colours.code"))
	} else if len(k.String("colours256.code")) != 0 {
		c.codeColour = numToAnsi("\033[38;5")(k.String("colours256.code"))
	}

	if len(k.String("colours.notify")) != 0 {
		c.notifyColour = hexToAnsi("\033[48;2")(k.String("colours.notify"))
	} else if len(k.String("colours256.notify")) != 0 {
		c.notifyColour = numToAnsi("\033[48;5")(k.String("colours256.notify"))
	}

}
