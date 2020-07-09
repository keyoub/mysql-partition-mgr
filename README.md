# MySQL Partition Manager

The `mspm` (MySQL Partition Manager) can be used to manage existing yearweek range paritioned MySQL tables.

```
NAME:
   spm - MySQL partition manager

USAGE:
   mspm [global options] command [command options] [arguments...]

VERSION:
   0.0.1

AUTHOR:
   Bardia Keyoumarsi <bardia@keyoumarsi.com>

COMMANDS:
   config-template    outputs a config file template
   validate-config    validate your configuration file
   status             status of current partitioned tables
   update-partitions  update database partitions based on given configuration
   help, h            Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help (default: false)
   --version, -v  print the version (default: false)
```

## config file template

```
./mspm config-template
{
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
```

### TODOs
* Support more range types like yearmonth

