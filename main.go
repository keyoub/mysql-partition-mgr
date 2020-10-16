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
	&cli.StringFlag{
		Name:     "database",
		Aliases:  []string{"db"},
		Usage:    "database name",
		Required: true,
		EnvVars:  []string{"DATABASE"},
	},
	&cli.StringFlag{
		Name:     "database-dsn",
		Aliases:  []string{"dsn"},
		Usage:    "database dsn",
		Required: true,
		EnvVars:  []string{"DATABASE_DSN"},
	},
}

var appCommands = []*cli.Command{
	{
		Name:   "template",
		Usage:  "outputs a config file template",
		Action: template,
	},
	{
		Name:   "yearweek",
		Usage:  "outputs the current yearweek value",
		Action: yearweek,
	},
	{
		Name:   "validate",
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
		Name:   "update",
		Usage:  "update database partitions based on given configuration",
		Action: update,
		Flags:  defaultFlags,
	},
}

func main() {
	app := &cli.App{
		Name:    "sqlpart",
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
