# Quake 3 Stats Server

## Building

Clone the code somewhere (`go get ..` works too):

```
go get github.com/bboozzoo/q3stats
```

Build with `make`. This will build both the server executable and
helper tools.

You can also run `make packed` to pack the binaries with UPX (make
sure that `upx` is available).

## Usage

### Server

First start the server process:

```
q3stats daemon
```

The daemon will open port `9090`, use your browser to access
http://localhost:9090

The deamon exports an API and provides a nice web view of your stats.

### Data Import

Match data from CPMA is imported using `q3simport` command. The
command will upload the match statistics to the server. If server is
running on another host/port, the address can be passed by adding `-t
<host>:<port>` in command line.

Example:

```
./q3simport matchdata/stats/2016/02/29/14_18_42.xml
```

The import tool is intentionally built as a separate program (and
lives under `cmd/q3simport` in the source tree). The for that is the
tool has very little (no) dependencies outside of Go's stdlib,
especially does not use `cgo` or link with external C libraries. This
allows for cross compiling `q3simport` to other architectures with
ease.

## License

MIT
