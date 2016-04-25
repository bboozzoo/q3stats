# Quake 3 Stats Server

## Building

```
go get github.com/bboozzoo/q3stats
```

## Usage

First start the daemon process:

```
q3stats daemon
```

The daemon will open port `9090`, use your browser to access
http://localhost:9090

Match data from CPMA is imported using `q3stats import` command. The
command will upload the match statistics to the server. If server is
runnin on another host/port, the address can be passed by adding `-t
<host>:<port>` in command line.

## License

MIT
