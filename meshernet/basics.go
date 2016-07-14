package meshernet

import (
	crand "crypto/rand"
	"fmt"
	"log"
	"net"
	"os"
	"os/user"
	"strconv"
	"time"

	"golang.org/x/crypto/nacl/box"
)

// InitMesher initialize internal meshernet stuff
func InitMesher() {
	neighborhood = make(map[string]Neighbor)

	userstruct, _ := user.Current()
	hostname, _ := os.Hostname()
	NodeName = userstruct.Username + "@" + hostname

	pubkey, privkey, err := box.GenerateKey(crand.Reader)
	CheckError(err)

	MySID = *pubkey
	MyPrivKey = privkey
	//fmt.Println(hex.EncodeToString(pubkey[:]))
	//fmt.Println(privkey)
}

// PrintPeers print all known peers
func PrintPeers() {
	for key, value := range neighborhood {
		fmt.Println(key, " : ", value.name, " : ", value.addr, " : ", time.Since(value.lastSeen))
	}
}

// ClearPeers forget all known peers
func ClearPeers() {
	for k := range neighborhood {
		delete(neighborhood, k)
	}
}

// CheckError Generic error checking, exit in case of error
func CheckError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(0)
	}
}

// ServerConn server listening socket, also used for sending
var ServerConn *net.UDPConn

// NodeName name of current node
var NodeName string

// MySID my public key
var MySID [32]byte

// MyPrivKey my private key
var MyPrivKey *[32]byte

// ListenerUDP start listening for incoming meshernet packets
func ListenerUDP(port int) chan bool {
	stop := make(chan bool)

	go func() {
		lport := ":" + strconv.Itoa(port)
		ServerAddr, err := net.ResolveUDPAddr("udp", lport)
		CheckError(err)
		ServerConnLocal, err := net.ListenUDP("udp", ServerAddr)
		CheckError(err)
		ServerConn = ServerConnLocal
		defer ServerConn.Close()

		buf := make([]byte, 1024)
		for {
			n, addr, err := ServerConn.ReadFromUDP(buf)
			//fmt.Println("Received ", string(buf[0:n]), " from ", addr)
			parsePkt(buf, n, addr)

			if err != nil {
				fmt.Println("Error: ", err)
			}
		}
	}()

	return stop
}

func sendPkt(dstStr string, port int, payload []byte) {
	lport := ":" + strconv.Itoa(port)

	dst, err := net.ResolveUDPAddr("udp", dstStr+lport)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := ServerConn.WriteToUDP(payload, dst); err != nil {
		log.Fatal(err)
	}
}

func sendPktToAddr(dst *net.UDPAddr, payload []byte) {
	if _, err := ServerConn.WriteToUDP(payload, dst); err != nil {
		log.Fatal(err)
	}
}
