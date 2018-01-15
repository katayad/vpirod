package main

import (
	"./Msg"
	"flag"
	"net"
	"strconv"
)

func main() {
	senderPtr := flag.Int("node", 0, "sender")
	actionPtr := flag.String("action", "send", "send/drop")
    dataPtr := flag.String("data", "msg", "data")
	destinationIdPtr := flag.Int("dest", 1, "Id")
	flag.Parse()

    sender := 40000 + *senderPtr

    ServerAddr, _ := net.ResolveUDPAddr("udp","127.0.0.1:" + strconv.Itoa(sender))
    LocalAddr, _ := net.ResolveUDPAddr("udp", "127.0.0.1:0")
    Conn, _ := net.DialUDP("udp", LocalAddr, ServerAddr)
    defer Conn.Close()

    msg := Msg.Msg{Type: *actionPtr, Dst: *destinationIdPtr, Src: *senderPtr, Data: *dataPtr, Accept: false}
	Conn.Write([]byte(msg.ToJson()))

}
