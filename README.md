# ping-go
ICMP echo request and response tool (known as ping) written in Go

# Getting dependencies
- For Windows, run get-ext.bat
- For Linux, run get.ext.sh

#### Requires 'git' for both Windows and Linux

Then, use the 'replace' function in your go.mod file to point the package to your new local directory. 
![replace github.com/jackpal/gateway => ./ext/gateway](https://i.imgur.com/wOBAvXD.png)

Why do you have to 'git clone' when you can just run 'go get'?
Because the dependency is outdated on 'go get' but is not outdated on 'git'. 

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
