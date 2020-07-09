package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/urfave/cli/v2"
)

const configTemplate = `{
	"database": "template",
	"database_dsn": "username:password@protocol(address)/dbname?param=value",
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
	Database    string  `json:"database"`
	DatabaseDSN string  `json:"database_dsn"`
	Tables      []Table `json:"tables"`
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

func config(c *cli.Context) (err error) {
	_, err = loadAndValidateConfig(c.String("config"))
	if err == nil {
		// don't leave the user hanging
		fmt.Fprintln(os.Stdout, "Configuration is valid!")
	}
	return
}

func loadAndValidateConfig(filePath string) (c *Config, err error) {
	configBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		err = errors.New("Unable to load config file: " + err.Error())
		return
	}

	err = json.Unmarshal(configBytes, &c)
	if err != nil {
		err = errors.New("Unable to parse config file: " + err.Error())
		return
	}

	db, err := connectDB(c.DatabaseDSN)
	if err != nil {
		err = errors.New("Failed to establish connection with SQL database: " + err.Error())
		return
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		err = errors.New("Failed to ping SQL database: " + err.Error())
		return
	}

	for _, table := range c.Tables {
		err = verifyTable(db, c.Database, table.Name, table.PartitionSchema)
		if err != nil {
			err = errors.New("Failed to verify table " + table.Name + ": " + err.Error())
			return
		}
	}

	return
}
