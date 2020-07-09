package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestUpdatePartitions(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	defer mockDB.Close()
	assert.NoError(t, err)

	rows := sqlmock.NewRows([]string{"PARTITION_NAME", "PARTITION_EXPRESSION", "PARTITION_DESCRIPTION", "TABLE_ROWS", "AVG_ROW_LENGTH", "DATA_LENGTH", "INDEX_LENGTH", "PARTITION_COMMENT"}).
		AddRow("p202019", "", "202019", 0, 0, 0, 0, "").
		AddRow("p202020", "", "202020", 0, 0, 0, 0, "").
		AddRow("p202021", "", "202021", 0, 0, 0, 0, "").
		AddRow("p202022", "", "202022", 0, 0, 0, 0, "").
		AddRow("p202023", "", "202023", 0, 0, 0, 0, "").
		AddRow("p202024", "", "202024", 0, 0, 0, 0, "").
		AddRow("p202025", "", "202025", 0, 0, 0, 0, "").
		AddRow("p202026", "", "202026", 0, 0, 0, 0, "")
	mock.ExpectQuery(`
		SELECT PARTITION_NAME, PARTITION_EXPRESSION, PARTITION_DESCRIPTION, TABLE_ROWS, AVG_ROW_LENGTH, DATA_LENGTH, INDEX_LENGTH, PARTITION_COMMENT
			FROM INFORMATION_SCHEMA.PARTITIONS
		WHERE TABLE_NAME = \? AND TABLE_SCHEMA = \?`).WillReturnRows(rows)

	mock.ExpectExec(`ALTER TABLE test_table ADD PARTITION \(PARTITION p202027 VALUES LESS THAN \(202027\)\)`).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(`ALTER TABLE test_table ADD PARTITION \(PARTITION p202028 VALUES LESS THAN \(202028\)\)`).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec("ALTER TABLE test_table DROP PARTITION p202019, p202020").WillReturnResult(sqlmock.NewResult(1, 1))

	err = updatePartitions(mockDB, "mocksql", Table{Name: "test_table", Retention: 4}, time.Date(2020, time.June, 19, 12, 0, 0, 0, time.UTC))
	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestUpdatePartitionsNoNewPartitions(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	defer mockDB.Close()
	assert.NoError(t, err)

	rows := sqlmock.NewRows([]string{"PARTITION_NAME", "PARTITION_EXPRESSION", "PARTITION_DESCRIPTION", "TABLE_ROWS", "AVG_ROW_LENGTH", "DATA_LENGTH", "INDEX_LENGTH", "PARTITION_COMMENT"}).
		AddRow("p202001", "", "202001", 0, 0, 0, 0, "").
		AddRow("p202002", "", "202002", 0, 0, 0, 0, "").
		AddRow("p202026", "", "202026", 0, 0, 0, 0, "").
		AddRow("p202027", "", "202027", 0, 0, 0, 0, "").
		AddRow("p202028", "", "202028", 0, 0, 0, 0, "")
	mock.ExpectQuery(`
		SELECT PARTITION_NAME, PARTITION_EXPRESSION, PARTITION_DESCRIPTION, TABLE_ROWS, AVG_ROW_LENGTH, DATA_LENGTH, INDEX_LENGTH, PARTITION_COMMENT
		  FROM INFORMATION_SCHEMA.PARTITIONS
		WHERE TABLE_NAME = \? AND TABLE_SCHEMA = \?`).WillReturnRows(rows)

	mock.ExpectExec("ALTER TABLE test_table DROP PARTITION p202001, p202002").WillReturnResult(sqlmock.NewResult(1, 1))

	err = updatePartitions(mockDB, "mocksql", Table{Name: "test_table", Retention: 4}, time.Date(2020, time.June, 19, 12, 0, 0, 0, time.UTC))
	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestUpdatePartitionsNoDropPartitions(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	defer mockDB.Close()
	assert.NoError(t, err)

	rows := sqlmock.NewRows([]string{"PARTITION_NAME", "PARTITION_EXPRESSION", "PARTITION_DESCRIPTION", "TABLE_ROWS", "AVG_ROW_LENGTH", "DATA_LENGTH", "INDEX_LENGTH", "PARTITION_COMMENT"}).
		AddRow("p202025", "", "202025", 0, 0, 0, 0, "").
		AddRow("p202026", "", "202026", 0, 0, 0, 0, "").
		AddRow("p202027", "", "202027", 0, 0, 0, 0, "")
	mock.ExpectQuery(`
		SELECT PARTITION_NAME, PARTITION_EXPRESSION, PARTITION_DESCRIPTION, TABLE_ROWS, AVG_ROW_LENGTH, DATA_LENGTH, INDEX_LENGTH, PARTITION_COMMENT
		  FROM INFORMATION_SCHEMA.PARTITIONS
		WHERE TABLE_NAME = \? AND TABLE_SCHEMA = ?`).WillReturnRows(rows)

	mock.ExpectExec(`ALTER TABLE test_table ADD PARTITION \(PARTITION p202028 VALUES LESS THAN \(202028\)\)`).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = updatePartitions(mockDB, "mocksql", Table{Name: "test_table", Retention: 4}, time.Date(2020, time.June, 19, 12, 0, 0, 0, time.UTC))
	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
