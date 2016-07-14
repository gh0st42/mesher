package main

import (
	"fmt"
	"log"
	"mesher/meshernet"
	"strings"
	"time"

	"github.com/chzyer/readline"
)

func schedule(what func(), delay time.Duration) chan bool {
	stop := make(chan bool)

	go func() {
		for {
			what()
			select {
			case <-time.After(delay):
			case <-stop:
				return
			}
		}
	}()

	return stop
}

var completer = readline.NewPrefixCompleter(
	readline.PcItem("/quit"),
	readline.PcItem("/msg"),
	readline.PcItem("/peers"),
	readline.PcItem("/clearpeers"),
	readline.PcItem("/help"))

func main() {
	meshernet.InitMesher()
	meshernet.ListenerUDP(8032)
	time.Sleep(500 * time.Millisecond)

	println("== MesherNet Mesh Chat ==")
	println(" Lars Baumgaertner (C)  2016")
	fmt.Println("\nServer started...\n")
	schedule(meshernet.SendAnnounce, 5000*time.Millisecond)

	println("Welcome ", meshernet.NodeName)
	rl, err := readline.NewEx(&readline.Config{
		UniqueEditLine: true,
		AutoComplete:   completer,
	})
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	username := meshernet.NodeName
	rl.ResetHistory()
	log.SetOutput(rl.Stderr())

	rl.SetPrompt(username + "> ")

	println(" Enter /help for command reference\n")

	for {
		ln := rl.Line()
		if ln.CanContinue() {
			continue
		} else if ln.CanBreak() {
			break
		}
		inputline := strings.TrimSpace(ln.Line)
		if inputline == "/quit" {
			println("\nQuitting...")
			break
		} else if inputline == "/peers" {
			println("\nKnown peers: ")
			meshernet.PrintPeers()
		} else if inputline == "/clearpeers" {
			println("Flushing peer list..")
			meshernet.ClearPeers()
			meshernet.PrintPeers()
		} else if strings.HasPrefix(inputline, "/msg") == true {
			msgfields := strings.SplitN(inputline, " ", 3)
			println(">> ", msgfields[1], " -> ", msgfields[2])
			if meshernet.SendPrivMsgByName(msgfields[1], msgfields[2]) == false {
				println("! No such name in peer list: ", msgfields[1])
			}
		} else if inputline == "/help" {
			println("\n-- MesherNet Commands --")
			println(" /quit - quit the program")
			println(" /peers - list all known peers")
			println(" /clearpeers - forget all known peers")
			println(" /msg <nodename> <message> - send a private message to node name\n")

		} else {
			//log.Println(username+":", inputline)
			meshernet.SendPubMsg(inputline)
		}
	}
	rl.Clean()
}
