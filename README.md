# ⚡ ️QuickTrace

## Intro

* This command line tool parse log file entries and build traces
* This program has been made for SimScale technical assessment

## How to use

 `Before anything make sure you have Go 1.12+ environment installed`.  
Link to docs : [https://golang.org/dl/](https://golang.org/dl/).  

```bash
$ git@github.com:cgstag/quicktrace.git
$ cd quicktrace
$ go build
$ go run quicktrace --input=stdin --output=stdout
```

## Test it

### Unit-tests

```bash
go test ./...
```

## Advanced Use

`--input=` : Choose input to read between standard input `stdin` or file `filename` (default **stdin**)

`--output=` : Choose output to write between standard output `stdout` or file `filename` (default **stdout**)

 ```bash
 ./quickstart --input=large-log.txt --output=large-quick-traces.txt
 ```

`--stats=true` : Print to standard error additional statistics and log messages (default **false**)

`--help` : Print use instructions


## Next Step

* Managing pending / orphans : In order to manage the traces without root span, one idea would be to add them in a memory queue with an expiration threshold. Two ideas come to mind :
    * The entry should batch-read a directory or be triggered through a pubsub topic
    * Redis as a queue looks like a good idea because of its builtin expiration