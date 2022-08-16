# Journal (jnl)


## Usage

```shell
jnl command --command-arg1 --command-arg2
```

Get command history

```shell
jnl getHistory
```

In CSV format: 
```shell
jnl getHistory --format=csv
```


## .jnl structure

The `.jnl` directory has directory for every command executed.
These directories are named after a hash computed from`command_executed:timestamp`.
Inside each one of those directories contains the following files:
- `stderr`: stderr output from the executed command
- `stdout`: stdout output from the executed command
- `info`: JSON file containing details about the execution



## TODO

- [ ] Save Output
