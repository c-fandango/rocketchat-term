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

var black = []string{
	"\033[38;5;0m",
}
var nothing = []string{
	"\033[0m",
}

const defaultCode = "\033[38;5;186m"
const defaultNotify = "\033[48;5;160m"

type configTemplate struct {
	userTextColours []string
	userBgColours   []string
	roomTextColours []string
	roomBgColours   []string
	codeColour      string
	notifyColour    string
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

func (c *configTemplate) loadConf(path string) {
	var k = koanf.New(".")

	if _, err := os.Stat(path); err == nil {
		k.Load(file.Provider(path), yaml.Parser())
	}

	// set colour defaults
	c.userTextColours = utils.MapperStr(defaultCols, numToAnsi("\033[38;5"))
	c.userBgColours = nothing
	c.roomTextColours = black
	c.roomBgColours = utils.MapperStr(defaultCols, numToAnsi("\033[48;5"))
	c.codeColour = defaultCode
	c.notifyColour = defaultNotify

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
