package meshernet

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"time"
)

func updateNeighborhood(hdr MeshPkt, sender *net.UDPAddr) {
	key := hex.EncodeToString(hdr.SenderID[:])

	entry := neighborhood[key]
	entry.identity = hdr.SenderID
	entry.lastSeen = time.Now()
	entry.addr = sender

	neighborhood[key] = entry
}
func parsePkt(pkt []byte, pktSize int, sender *net.UDPAddr) {
	var hdr MeshPkt
	buffer := bytes.NewBuffer(pkt)
	err := binary.Read(buffer, binary.BigEndian, &hdr)
	if err != nil {
		fmt.Println("binary.Read failed:", err)
	}
	updateNeighborhood(hdr, sender)
	//println("RECEIVED: ", hdr.PktPayload)
	switch hdr.PktPayload {
	case TypeAnnounce:
		var packet AnnouncePkt
		err := binary.Read(buffer, binary.BigEndian, &packet)
		if err != nil {
			fmt.Println("binary.Read failed:", err)
		}
		handleAnnounce(hdr, packet)

	case TypePubMsg:
		var packet PubMsgPkt
		err := binary.Read(buffer, binary.BigEndian, &packet)
		if err != nil {
			fmt.Println("binary.Read failed:", err)
		}
		msgtext := string(buffer.Next(int(packet.MsgLen)))
		handlePubMsg(hdr, packet, msgtext)

	case TypePrivMsg:
		var packet PrivMsgPkt
		err := binary.Read(buffer, binary.BigEndian, &packet)
		if err != nil {
			fmt.Println("binary.Read failed:", err)
		}
		msgtext := string(buffer.Next(int(packet.MsgLen)))
		handlePrivMsg(hdr, packet, msgtext)

	case TypeLookupRep:
		var packet LookupRepPkt
		err := binary.Read(buffer, binary.BigEndian, &packet)
		if err != nil {
			fmt.Println("binary.Read failed:", err)
		}
		nodename := string(buffer.Next(int(packet.NameLen)))
		handleLookupRep(hdr, packet, nodename)

	case TypeLookupReq:
		var packet LookupReqPkt
		err := binary.Read(buffer, binary.BigEndian, &packet)
		if err != nil {
			fmt.Println("binary.Read failed:", err)
		}
		handleLookupReq(hdr, packet)

	default:
		fmt.Printf("unknown packet type: %d\n", hdr.PktPayload)
	}
	//fmt.Println(hdr)
}

func handleLookupReq(hdr MeshPkt, pkt LookupReqPkt) {
	SendLookupReply(ReplyPeer(hdr))
}

func handleLookupRep(hdr MeshPkt, pkt LookupRepPkt, nodename string) {
	key := hex.EncodeToString(hdr.SenderID[:])

	entry := neighborhood[key]
	entry.name = nodename
	neighborhood[key] = entry
}

func handleAnnounce(hdr MeshPkt, pkt AnnouncePkt) {
	key := hex.EncodeToString(hdr.SenderID[:])

	entry := neighborhood[key]
	if entry.name == "" {
		SendLookupRequest(entry)
	}
}
func handlePubMsg(hdr MeshPkt, pkt PubMsgPkt, msgtext string) {
	peer := ReplyPeer(hdr)
	log.Println(peer.name, ": ", msgtext)
}
func handlePrivMsg(hdr MeshPkt, pkt PrivMsgPkt, msgtext string) {
	ciphertext, _ := hex.DecodeString(msgtext)

	cleartext, _ := Decrypt(MyPrivKey, &hdr.SenderID, ciphertext)

	peer := ReplyPeer(hdr)
	log.Println("PRIVATE (", peer.name, "): ", string(cleartext))
}
