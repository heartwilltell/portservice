# 🚢 ports service

Represents a simple service which receives the path to files that contains information about ports.

## Build

To build the service binary run the following command:

```shell
GOOS=linux GOARCH=amd64 go build -ldflags="-X main.Branch=$BRANCH -X main.Commit=$COMMIT" -o service ./cmd
```

_Don't forget to set proper GOOS and GOARCH values for target platform and architecture._

## Test

To run tests run the following command:

```shell
go test -race -cover ./...
```

## Lint

```shell
golangci-lint run ./...
```

_The golangci-lint executable should be installed on the system._

## Usage

```text
NAME:
   app - root command of the app

USAGE:
   app [global options] command [command options] [arguments...]

COMMANDS:
   run      runs worker to process files with port information
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help
```

The `run` command usage:

```text
NAME:
   app run - runs worker to process files with port information

USAGE:
   app run [command options] [arguments...]

OPTIONS:
   --worker-conc-limit value    defines worker concurrency limit (default: 8) [$WORKER_CONC_LIMIT]
   --worker-proc-timeout value  defines worker processing timeout (default: 5m0s) [$WORKER_PROC_TIMEOUT]
   --ports-file-path value      sets the path to json file with ports [$PORTS_FILE_PATH]
   --help, -h                   show help
```