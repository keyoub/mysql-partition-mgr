package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/urfave/cli/v2"
)

func update(c *cli.Context) error {
	conf, err := loadAndValidateConfig(c.String("config"))
	if err != nil {
		return err
	}

	db, err := connectDB(conf.DatabaseDSN)
	if err != nil {
		return errors.New("Failed to establish connection with SQL database: " + err.Error())
	}
	defer db.Close()

	for _, table := range conf.Tables {
		err = updatePartitions(db, conf.Database, table, time.Now())
		if err != nil {
			fmt.Fprintln(os.Stdout, err.Error())
		}
	}

	printStatus(db, conf)

	return nil
}
