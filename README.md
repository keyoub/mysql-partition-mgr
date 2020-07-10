# MySQL Partition Manager

The `mspm` (MySQL Partition Manager) can be used to manage existing yearweek range paritioned MySQL tables.

```
NAME:
   spm - MySQL partition manager

USAGE:
   mspm [global options] command [command options] [arguments...]

VERSION:
   1.0.0

AUTHOR:
   Bardia Keyoumarsi <bardia@keyoumarsi.com>

COMMANDS:
   config-template    outputs a config file template
   yearweek           outputs the current yearweek value
   validate-config    validate your configuration file
   status             status of current partitioned tables
   update-partitions  update database partitions based on given configuration
   help, h            Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help (default: false)
   --version, -v  print the version (default: false)
```

## Installation

```bash
git clone github.com/qeubar/mysql-partition-mgr
cd mysql-partition-mgr/
go install
```

## Example

create table:
```MySQL
CREATE TABLE `logs` (
  `id` BIGINT(21) NOT NULL AUTO_INCREMENT,
  `yearweek` INT(11) NOT NULL,
  `user_id` BIGINT(21) NOT NULL,
  `message` TEXT NOT NULL,
  `created_at` BIGINT(21) NOT NULL,
  PRIMARY KEY (`id`, `yearweek`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
PARTITION BY RANGE(yearweek) (
      PARTITION p202024 VALUES LESS THAN (202024),
      PARTITION p202025 VALUES LESS THAN (202025)
);
```

setup config:
```bash
cat > myconfig.json
{
        "database": "mydatabase",
        "database_dsn": "username:password@protocol(address)/mydatabase",
        "tables": [
                {
                        "name": "logs",
                        "partition_schema": "yearweek",
                        "retention": 5,
                        "max_future_partitions": 5
                }
        ]
}
```

run update:
```bash
./mspm update-partitions -c myconfig.json
Partitions to add to the logs table [202029 202030 202031 202032]
+-------+----------------+----------------------+-----------------------+----------------+-----------------------+-----------------+-------------------+---------+
| TABLE | PARTITION NAME | PARTITION EXPRESSION | PARTITION DESCRIPTION | NUMBER OF ROWS | AVERAGE ROW SIZE (MB) | INDEX SIZE (MB) | STORAGE SIZE (MB) | COMMENT |
+-------+----------------+----------------------+-----------------------+----------------+-----------------------+-----------------+-------------------+---------+
| logs  | p202024        | yearweek             |                202024 |              0 |                     0 |               0 |          0.015625 |         |
| logs  | p202025        | yearweek             |                202025 |              0 |                     0 |               0 |          0.015625 |         |
| logs  | p202029        | yearweek             |                202029 |              0 |                     0 |               0 |          0.015625 |         |
| logs  | p202030        | yearweek             |                202030 |              0 |                     0 |               0 |          0.015625 |         |
| logs  | p202031        | yearweek             |                202031 |              0 |                     0 |               0 |          0.015625 |         |
| logs  | p202032        | yearweek             |                202032 |              0 |                     0 |               0 |          0.015625 |         |
+-------+----------------+----------------------+-----------------------+----------------+-----------------------+-----------------+-------------------+---------+
```

### TODOs
* Support more range types like yearmonth
* Publish binary to select distros
