# ping-go
ICMP echo request and response tool (known as ping) written in Go

# Usage
### For Linux
(building from source)
```console
foo@bar:~$ go build
foo@bar:~$ ./ping-go [-c count] destination
```
if you get `socket: operation not permitted`, try running as root (or sudo)

### For Windows
(building from source)
```bat
go build
ping-go.exe [-c count] destination
```
