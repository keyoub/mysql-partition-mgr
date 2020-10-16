package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/urfave/cli/v2"
)

const configTemplate = `{
	"tables": [
		{
			"name": "x",
			"partition_schema": "yearweek",
			"retention": 5,
			"max_future_partitions": 5
		},
		{
			"name": "y",
			"partition_schema": "yearweek",
			"retention": 2,
			"max_future_partitions": 1
		}
	]
}
`

type Config struct {
	Database    string
	DatabaseDSN string
	Tables      []Table `json:"tables"`

	db *sql.DB
}

type Table struct {
	Name                string `json:"name"`
	PartitionSchema     string `json:"partition_schema"`
	Retention           int    `json:"retention"`
	MaxFuturePartitions int    `json:"max_future_partitions"`
}

func template(_ *cli.Context) error {
	fmt.Fprintln(os.Stdout, configTemplate)
	return nil
}

func config(ctx *cli.Context) (err error) {
	_, err = loadAndValidateConfig(ctx)
	if err == nil {
		// don't leave the user hanging
		fmt.Fprintln(os.Stdout, "Configuration is valid!")
	}
	return
}

func loadAndValidateConfig(ctx *cli.Context) (*Config, error) {
	c := Config{
		Database:    ctx.String("database"),
		DatabaseDSN: ctx.String("database-dsn"),
	}

	configBytes, err := ioutil.ReadFile(ctx.String("config"))
	if err != nil {
		err = errors.New("Unable to load config file: " + err.Error())
		return nil, err
	}

	err = json.Unmarshal(configBytes, &c)
	if err != nil {
		err = errors.New("Unable to parse config file: " + err.Error())
		return nil, err
	}

	c.db, err = connectDB(c.DatabaseDSN)
	if err != nil {
		err = errors.New("Failed to establish connection with SQL database: " + err.Error())
		return nil, err
	}

	err = c.db.Ping()
	if err != nil {
		err = errors.New("Failed to ping SQL database: " + err.Error())
		return nil, err
	}

	for _, table := range c.Tables {
		err = verifyTable(c.db, c.Database, table.Name, table.PartitionSchema)
		if err != nil {
			err = errors.New("Failed to verify table " + table.Name + ": " + err.Error())
			return nil, err
		}
	}

	return &c, nil
}
