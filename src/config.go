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

const defaultCode = "186"
const defaultNotify = "160"

type configTemplate struct {
	textColours  []string
	bgColours    []string
	codeColour   string
	notifyColour string
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

	if len(k.Strings("colours.text")) != 0 {
		c.textColours = utils.MapperStr(k.Strings("colours.text"), hexToAnsi("\033[38;2"))
	} else if len(k.Strings("colours256.text")) != 0 {
		c.textColours = utils.MapperStr(k.Strings("colours256.text"), numToAnsi("\033[38;5"))
	} else {
		c.textColours = utils.MapperStr(defaultCols, numToAnsi("\033[38;5"))
	}

	if len(k.Strings("colours.highlights")) != 0 {
		c.bgColours = utils.MapperStr(k.Strings("colours.highlights"), hexToAnsi("\033[48;2"))
	} else if len(k.Strings("colours256.highlights")) != 0 {
		c.bgColours = utils.MapperStr(k.Strings("colours256.highlights"), numToAnsi("\033[48;5"))
	} else {
		c.bgColours = utils.MapperStr(defaultCols, numToAnsi("\033[48;5"))
	}

	if len(k.String("colours.code")) != 0 {
		c.codeColour = hexToAnsi("\033[38;2")(k.String("colours.code"))
	} else if len(k.String("colours256.code")) != 0 {
		c.codeColour = numToAnsi("\033[38;5")(k.String("colours256.code"))
	} else {
		c.codeColour = numToAnsi("\033[38;5")(defaultCode)
	}

	if len(k.String("colours.notify")) != 0 {
		c.notifyColour = hexToAnsi("\033[48;2")(k.String("colours.notify"))
	} else if len(k.String("colours256.notify")) != 0 {
		c.notifyColour = numToAnsi("\033[48;5")(k.String("colours256.notify"))
	} else {
		c.notifyColour = numToAnsi("\033[48;5")(defaultNotify)
	}
}
