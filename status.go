package main

import (
	"database/sql"
	"errors"
	"os"

	"github.com/jedib0t/go-pretty/table"
	"github.com/urfave/cli/v2"
)

func printStatus(db *sql.DB, conf *Config) error {
	status := make(map[string][]Partition)
	for _, table := range conf.Tables {
		p, err := getCurrPartitions(db, conf.Database, table.Name)
		if err != nil {
			return err
		}
		status[table.Name] = p
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(
		table.Row{
			"Table",
			"Partition Name",
			"Partition Expression",
			"Partition Description",
			"Number of Rows",
			"Average Row Size (MB)",
			"Index Size (MB)",
			"Storage Size (MB)",
			"Comment",
		},
	)
	for tableName, partitions := range status {
		for _, p := range partitions {
			t.AppendRows(
				[]table.Row{{
					tableName,
					p.Name,
					p.Expression,
					p.Description,
					p.TableRows,
					float32(p.AvgRowLength) / (1024 * 1024),
					float32(p.IndexLength) / (1024 * 1024),
					float32(p.DataLength) / (1024 * 1024),
					p.PartitionComment,
				}},
			)
		}
	}
	t.Render()

	return nil
}

func status(c *cli.Context) (err error) {
	conf, err := loadAndValidateConfig(c.String("config"))
	if err != nil {
		return
	}

	db, err := connectDB(conf.DatabaseDSN)
	if err != nil {
		err = errors.New("Failed to establish connection with SQL database: " + err.Error())
		return
	}
	defer db.Close()

	return printStatus(db, conf)
}
