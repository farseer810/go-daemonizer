# go-daemonizer

go-daemonizer is tool for daemonizing go server application. The underlying OS api is unix/linux fork, which means that this tool only works on unix-like systems(not on Windows). 

This project is created for the purpose of upgrading server application gracefully, though only achieve part of it. I was planning to add signal control to make it enable to quit, restart, upgrade gracefully, as nginx does. I then choose the ability of later extension instead of make this project specific purpose. If you really do want the other funtionalities, all you have to do is to add signal control.

## Getting Started

To install this, run
~~~
go get github.com/farseer810/go-daemonizer
~~~

Here's an all-in-one example:
~~~go
package main

import (
	"github.com/farseer810/go-daemonizer"
	"log"
	"net"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusOK)
		res.Write([]byte("hello, world"))
	})

	var listener net.Listener
	var address string = ":3000"
	var server *http.Server = &http.Server{Addr: address, Handler: nil}
	var err error

	if !daemonizer.IsDaemon {
		listener, err = net.Listen("tcp", address)
		if err != nil {
			log.Fatalln(err)
		}
		tcpListener, _ := listener.(*net.TCPListener)
		listenerFile, _ := tcpListener.File()
		daemonizer.AddSharedFileDescriptor(listenerFile.Fd())

		pid, err := daemonizer.Daemonize()
		if err != nil {
			log.Fatalln(err)
		}
		_ = pid
		return
	} else {
		listener, err = net.FileListener(os.NewFile(daemonizer.SharedFileDescriptors()[0], "arbitrary name"))
		if err != nil {
			log.Fatalln(err)
		}
	}
	log.Fatalln(server.Serve(listener))
}
~~~

Prototypes provided here:
~~~go
var (
	IsDaemon bool
)

func AddSharedFileDescriptor(fd uintptr)
func SharedFileDescriptors() []uintptr
func Daemonize() (pid int, err error)
~~~

## Shared File Descriptor

To add file descriptors which will be inherited by subprocesses, simply call AddSharedFileDescriptor:
~~~go
daemonizer.AddSharedFileDescriptor(listenerFile.Fd())
~~~

To get added file descriptors, call SharedFileDescriptors:
~~~go
daemonizer.SharedFileDescriptors()
~~~

## Daemonization
Simply run: 
~~~go
daemonizer.Daemonize()
~~~
which returns pid of the child process and a error(if any).

The variable daemonizer.IsDaemon(type bool) tells you whether this process is a child process.
