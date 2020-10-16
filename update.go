package main

import (
	"fmt"
	"os"
	"time"

	"github.com/urfave/cli/v2"
)

func update(c *cli.Context) error {
	conf, err := loadAndValidateConfig(c)
	if err != nil {
		return err
	}
	defer conf.db.Close()

	for _, table := range conf.Tables {
		err = updatePartitions(conf.db, conf.Database, table, time.Now())
		if err != nil {
			fmt.Fprintln(os.Stdout, err.Error())
		}
	}

	printStatus(conf)

	return nil
}
