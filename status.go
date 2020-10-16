package main

import (
	"fmt"
	"os"
	"time"

	"github.com/jedib0t/go-pretty/table"
	"github.com/urfave/cli/v2"
)

func printStatus(conf *Config) error {
	status := make(map[string][]Partition)
	for _, table := range conf.Tables {
		p, err := getCurrPartitions(conf.db, conf.Database, table.Name)
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
	conf, err := loadAndValidateConfig(c)
	if err != nil {
		return
	}

	return printStatus(conf)
}

func yearweek(_ *cli.Context) error {
	now := time.Now()
	yw := determineYearWeek(now)
	fmt.Fprintf(os.Stdout, "%d is the yearweek for %s\n", yw, now.Format(time.UnixDate))
	return nil
}
