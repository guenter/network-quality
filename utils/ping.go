package utils

import (
	"time"
	"os"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"log"
	"net"
	"fmt"
	"errors"
)

const ProtocolICMP = 1
const pingTimeout = 5 * time.Second

func Ping(networkInterface string, targetHost string) (time.Duration, error) {
	c, err := icmp.ListenPacket("udp4", "0.0.0.0")
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	c.SetDeadline(time.Now().Add(pingTimeout))

	wm := icmp.Message{
		Type: ipv4.ICMPTypeEcho, Code: 0,
		Body: &icmp.Echo{
			ID: os.Getpid() & 0xffff, Seq: 1,
			Data: []byte("HI!"),
		},
	}
	wb, err := wm.Marshal(nil)
	if err != nil {
		log.Fatal(err)
	}

	// START PING
	startTime := time.Now()

	if _, err := c.WriteTo(wb, &net.UDPAddr{IP: net.ParseIP(targetHost), Zone: networkInterface}); err != nil {
		log.Fatal(err)
	}

	rb := make([]byte, 1500)
	n, peer, err := c.ReadFrom(rb)
	if err != nil {
		log.Fatal(err)
	}

	// END PING
	endTime := time.Now()

	rm, err := icmp.ParseMessage(ProtocolICMP, rb[:n])
	if err != nil {
		log.Fatal(err)
	}

	switch rm.Type {
	case ipv4.ICMPTypeEchoReply:
		return endTime.Sub(startTime), nil
	default:
		return time.Duration(1<<63 - 1), errors.New(fmt.Sprintf("got %+v from %v; want echo reply", rm, peer))
	}
}