# mesher

A simple local IPv4 based mesh chat.

Features:
- peer-to-peer local group chat
- end-to-end encrypted private messages
- simple, extensible protocol

**WARNING THIS SOFTWARE IS NOT MEANT TO BE USED FOR ANYTHING SERIOUS!**

It was written for an IT security lecture at university to have something for the students to practice protocol reverse engineering and the development of wireshark dissectors!

## Installation, dependencies and running

Get all deps:
`go get ./...`

Directly running *mesher* :
`go run mesher.go`

or

```
$ go build mesher
$ ./mesher
```

or 

```
$ go install mesher
$ $GOPATH/bin/mesher
```

**Don't forget to set you GOPATH**