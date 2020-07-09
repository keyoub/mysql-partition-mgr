package main

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Partition struct {
	Name             string // p202030
	Expression       string // yearweek
	Description      int    // 202030
	TableRows        int    // 0
	AvgRowLength     int    // 0
	DataLength       int64  // 16384
	IndexLength      int64  // 16384
	PartitionComment string //
}

// determineYearWeek returns a YYYYWW integer.
func determineYearWeek(t time.Time) int {
	year, week := t.ISOWeek()

	return (year*100 + week)
}

// determineYearMonth returns a YYYYMM integer.
func determineYearMonth(t time.Time) int {
	year := t.Year()
	month := int(t.Month())

	return (year*100 + month)
}

// connectDB initiates a sql.DB instance, the caller is responsible for calling sql.DB.Close()
func connectDB(dsn string) (dbh *sql.DB, err error) {
	dbh, err = sql.Open("mysql", dsn)
	if err != nil {
		return
	}

	// Set really low limits, this application is only meant to do quick serialized SQL queries
	dbh.SetMaxOpenConns(1)
	dbh.SetConnMaxLifetime(time.Second)

	return
}

func verifyTable(dbh *sql.DB, dbName, tableName, partitionSchema string) error {
	var q = `
		SELECT PARTITION_DESCRIPTION, PARTITION_EXPRESSION, PARTITION_METHOD
		  FROM INFORMATION_SCHEMA.PARTITIONS
		WHERE TABLE_NAME = ? AND TABLE_SCHEMA = ? LIMIT 1`

	var partitionName, currPartitionSchema, partitionMethod string
	err := dbh.QueryRow(
		q,
		tableName,
		dbName,
	).Scan(
		&partitionName,
		&currPartitionSchema,
		&partitionMethod,
	)
	switch err {
	case nil:
		break
	case sql.ErrNoRows:
		return errors.New("Table does not exist")
	default:
		return err
	}

	if partitionName == "" {
		return errors.New("Partitioning is not configured")
	}

	if strings.ToLower(partitionSchema) != strings.ToLower(currPartitionSchema) {
		return errors.New("Table's actual partition schema is " + currPartitionSchema)
	}

	if strings.ToLower(partitionMethod) != "range" {
		return errors.New(`Table must have "RANGE" partitions`)
	}

	return nil
}

func getCurrPartitions(dbh *sql.DB, dbName, tableName string) (currPartitions []Partition, err error) {
	var q = `
		SELECT
			PARTITION_NAME,
			PARTITION_EXPRESSION,
			PARTITION_DESCRIPTION,
			TABLE_ROWS,
			AVG_ROW_LENGTH,
			DATA_LENGTH,
			INDEX_LENGTH,
			PARTITION_COMMENT
		FROM INFORMATION_SCHEMA.PARTITIONS
		WHERE TABLE_NAME = ? AND TABLE_SCHEMA = ?`

	rows, err := dbh.Query(q, tableName, dbName)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var p Partition
		err = rows.Scan(
			&p.Name,
			&p.Expression,
			&p.Description,
			&p.TableRows,
			&p.AvgRowLength,
			&p.DataLength,
			&p.IndexLength,
			&p.PartitionComment,
		)
		if err != nil {
			return
		}
		currPartitions = append(currPartitions, p)
	}

	return
}

func updatePartitions(dbh *sql.DB, dbName string, table Table, now time.Time) error {
	var partitionsToAdd []int
	var partitionsToDrop []string
	var currentPartitions []Partition

	oldestYearweekToDrop := determineYearWeek(now.AddDate(0, 0, table.Retention*-7))

	mfp := table.MaxFuturePartitions
	if mfp < 1 {
		mfp = 4 // default to 3 future partitions
	}
	for i := 1; i < mfp; i++ {
		partitionsToAdd = append(partitionsToAdd, determineYearWeek(now.AddDate(0, 0, i*7)))
	}

	currentPartitions, err := getCurrPartitions(dbh, dbName, table.Name)
	if err != nil {
		return err
	}

	for _, partition := range currentPartitions {
		partitionYearweek := partition.Description

		if partitionYearweek <= oldestYearweekToDrop {
			partitionsToDrop = append(partitionsToDrop, partition.Name)
		}

		for i, p := range partitionsToAdd {
			if partitionYearweek == p {
				// Partition already exists, remove it from ToAdd list
				partitionsToAdd = append(partitionsToAdd[:i], partitionsToAdd[i+1:]...)
			}
		}
	}

	if len(partitionsToAdd) > 0 {
		fmt.Fprintf(os.Stdout, "Partitions to add to the %s table %v\n", table.Name, partitionsToAdd)
		for _, partition := range partitionsToAdd {
			qAddPartition := fmt.Sprintf("ALTER TABLE %s ADD PARTITION (PARTITION p%d VALUES LESS THAN (%d))",
				table.Name, partition, partition)

			_, err = dbh.Exec(qAddPartition)
			if err != nil {
				return fmt.Errorf("Failed to add new partition %d: %s", partition, err)
			}
		}
	}

	if len(partitionsToDrop) > 0 {
		fmt.Fprintf(os.Stdout, "Partitions to drop from the %s table %v\n", table.Name, partitionsToDrop)
		qDropPartitions := fmt.Sprintf("ALTER TABLE %s DROP PARTITION %s",
			table.Name, strings.Join(partitionsToDrop, ", "))

		_, err = dbh.Exec(qDropPartitions)
		if err != nil {
			return fmt.Errorf("Failed to drop old partitions %v: %s", partitionsToDrop, err)
		}
	}

	return nil
}
