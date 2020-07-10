package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

var defaultFlags = []cli.Flag{
	&cli.StringFlag{
		Name:     "config",
		Aliases:  []string{"c"},
		Usage:    "configuration file",
		Required: true,
	},
}

var appCommands = []*cli.Command{
	{
		Name:   "config-template",
		Usage:  "outputs a config file template",
		Action: template,
	},
	{
		Name:   "yearweek",
		Usage:  "outputs the current yearweek value",
		Action: yearweek,
	},
	{
		Name:   "validate-config",
		Usage:  "validate your configuration file",
		Action: config,
		Flags:  defaultFlags,
	},
	{
		Name:   "status",
		Usage:  "status of current partitioned tables",
		Action: status,
		Flags:  defaultFlags,
	},
	{
		Name:   "update-partitions",
		Usage:  "update database partitions based on given configuration",
		Action: update,
		Flags:  defaultFlags,
	},
}

func main() {
	app := &cli.App{
		Name:    "spm",
		Usage:   "MySQL partition manager",
		Version: "1.0.0",
		Authors: []*cli.Author{
			{Name: "Bardia Keyoumarsi", Email: "bardia@keyoumarsi.com"},
		},
		Commands: appCommands,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
