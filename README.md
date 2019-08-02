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

2013-10-23T10:12:37.634Z 2013-10-23T10:12:37.946Z rjopvy3w service8 null->qa5oegbm

```
Setting stdin as input will wait for the next keyboard input.

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


## Next Steps

* Pending entries management
   * One idea would be to store in memory the pending entries, and try every second or so to find their parents in a separate goroutine. 
   * Upon trigger of that goroutine, it could compare the timestamp of that entry with the current timestamp (which is already stored in a global variable of the quicktrace package), and add make it an orphan in case of expiration. The expiration could be parametrically defined via command line flags
