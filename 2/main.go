package main

import "fmt"
import "time"
import "net"
import "flag"
import "./msg"
import "os"
import "strconv"

func main() {
    nodesCntPtr := flag.Int("n", 0, "nodes cnt")
    withholdingPtr := flag.Int("t", 0, "withholding time")
    flag.Parse()

    doneChan := make(chan int)
    cnt := 0
	for i := 0; i < *nodesCntPtr; i++ {
		go start(i, *nodesCntPtr, *withholdingPtr, 0, doneChan)
        cnt++;
	}
	for i := 0; i < *nodesCntPtr; i++ {
		<-doneChan
	}

    fmt.Println(cnt)

}

func start(nodeId int, nodesCnt int, withholding int, leaderId int, doneChan chan int) {
    msgPort, servicePort := 30000, 40000
    next := strconv.Itoa(30000 + (nodeId + 1) % nodesCnt)

    mainChan := make(chan Msg.Msg)
    serviceChan := make(chan Msg.Msg)
    go startServer(strconv.Itoa(msgPort + nodeId), mainChan)
    go startServer(strconv.Itoa(servicePort + nodeId), serviceChan)


    ServerAddr,err := net.ResolveUDPAddr("udp","127.0.0.1:" + next)
    CheckError(err, 1)
    LocalAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
    CheckError(err, 2)
    nextConn, err := net.DialUDP("udp", LocalAddr, ServerAddr)
    CheckError(err, 3)
    defer nextConn.Close()

    emptyMsg := Msg.Msg{Type: "", Dst: -1, Data: "", Accept: false, Src: 0, LastTouched: leaderId}
    msgToSend := emptyMsg
    dropTocken := false

    timer := time.NewTimer(time.Second * 4)
    for {
        select {
        case msg := <- mainChan:
            lastTouched := msg.LastTouched
            msg.LastTouched = nodeId
            emptyMsg.LastTouched = nodeId
            msgToSend.LastTouched = nodeId
            time.Sleep(time.Second * time.Duration(withholding))
            if dropTocken {
                fmt.Println("node ", nodeId, ": OOPS!")
                dropTocken = false
            } else if msg.Empty() {
                fmt.Println(
                    "node ", nodeId, ": recieved token from node ", lastTouched,
                    ", sending tocken to node ", (nodeId + 1) % nodesCnt)
                if msgToSend.Empty() {
                    sendUdp(msg, nextConn, 1)
                } else {
                    sendUdp(msgToSend, nextConn, 2)
                }
            } else {
                if msg.Dst == nodeId {
                    if msg.Accept {
                        fmt.Println(
                            "node ", nodeId, ": recieved token from node ", lastTouched,
                            " with delivery confirmation from node ", msg.Src,
                            ", sending tocken to node ", (nodeId + 1) % nodesCnt)
                        msgToSend = emptyMsg
                        sendUdp(emptyMsg, nextConn, 3)
                    } else {
                        fmt.Println(
                            "node ", nodeId, ": recieved token from node ", lastTouched,
                            " with data from node \n", msg.Src, "    (data='", msg.Data, "')",
                            ", sending tocken to node ", (nodeId + 1) % nodesCnt)
                        msg = Msg.Msg{Type: "", Dst: msg.Src, Data: "", Accept: true, Src: nodeId}
                        sendUdp(msg, nextConn, 4)
                    }
                } else {
                    fmt.Println(
                        "node ", nodeId, ": recieved token from node ", lastTouched,
                        ", sending tocken to node ", (nodeId + 1) % nodesCnt)
                    sendUdp(msg, nextConn, 5)
                }
            }
            timer = time.NewTimer(time.Second * time.Duration(withholding * nodesCnt + nodesCnt / 4) )
        case msg := <- serviceChan:
            fmt.Println("node ", nodeId, ": recieved service message\n    ", msg.ToJson())
            if msg.Type == "send" {
                msgToSend = msg
            } else if msg.Type == "drop" {
                dropTocken = true
            }

        case <- timer.C:
            if nodeId == leaderId {
                fmt.Println("node ", nodeId, ": Tocken fucked up, starting new one")
                sendUdp(emptyMsg, nextConn, 5)
            }
            timer = time.NewTimer(time.Second * time.Duration(withholding * nodesCnt + nodesCnt / 4) )
        }
    }

    doneChan <- 1
}

func startServer(port string, msgChan chan Msg.Msg) {
    ServerAddr, err := net.ResolveUDPAddr("udp",":" + port)
    CheckError(err, 4)

    ServerConn, err := net.ListenUDP("udp", ServerAddr)
    CheckError(err, 5)

    defer ServerConn.Close()

    buf := make([]byte, 1024)
    for {
        n, _, err := ServerConn.ReadFromUDP(buf)
        CheckError(err, 6)
        msgChan <- Msg.New(buf[0:n])
    }
}

func sendUdp(msg Msg.Msg, ServerConn *net.UDPConn, flag int) {
    buf := []byte(msg.ToJson())
    _,err := ServerConn.Write(buf)
    CheckError(err, 100 + flag)
}

func CheckError(err error, flag int) {
    if err != nil {
        fmt.Println("Error: " , err, "   flag", flag)
        os.Exit(0)
    }
}
