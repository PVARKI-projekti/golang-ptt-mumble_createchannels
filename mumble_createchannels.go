package main

import (
	"os"

	"github.com/romnn/flags4urfavecli/flags"
	"github.com/romnn/flags4urfavecli/values"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

// Rev is set on build time to the git HEAD
var Rev = ""

// Version is incremented using bump2version
const Version = "0.0.2"

// Shout returns the input message with an exclamation mark
func Shout(s string) string {
	return s + "!"
}

func main() {
	app := &cli.App{
		Name:  "mumble_createchannels",
		Usage: "",
		Flags: []cli.Flag{
			&cli.GenericFlag{
				Name: "format",
				Value: &values.EnumValue{
					Enum:    []string{"json", "xml", "csv"},
					Default: "xml",
				},
				EnvVars: []string{"FILEFORMAT"},
				Usage:   "input file format",
			},
			&flags.LogLevelFlag,
		},
		Action: func(ctx *cli.Context) error {
			if level, err := log.ParseLevel(ctx.String("log")); err == nil {
				log.SetLevel(level)
			}
			log.Infof("Format is: %s", ctx.String("format"))
			return nil
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
