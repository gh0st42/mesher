package meshernet

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	//"fmt"
	"net"
)

// SendAnnounce sends announce message to broadcast address
func SendAnnounce() {
	var pkt MeshPkt
	var annc AnnouncePkt

	pkt.PktPayload = TypeAnnounce
	pkt.SenderID = MySID

	annc.NumServices = 3
	annc.Services = [32]int8{ServiceLookup, ServicePubChat, ServicePrivMsg}

	var buffer bytes.Buffer

	binary.Write(&buffer, binary.BigEndian, &pkt)
	binary.Write(&buffer, binary.BigEndian, &annc)

	sendPkt("255.255.255.255", 8032, buffer.Bytes())
}

// SendPrivMsgByName helper function to send private message to nodename
func SendPrivMsgByName(receiver string, message string) bool {
	for _, value := range neighborhood {
		if value.name == receiver {
			SendPrivMsg(message, value.identity)
			return true
		}
	}
	return false
}

// SendPrivMsg send private message to specified peer
func SendPrivMsg(message string, receiverID [32]byte) {
	key := hex.EncodeToString(receiverID[:])

	ciphertext, _ := Encrypt(MyPrivKey, &receiverID, []byte(message))

	cryptomsg := hex.EncodeToString(ciphertext)

	pkt := MeshPkt{TypePrivMsg, MySID}
	payload := PrivMsgPkt{receiverID, int32(len(cryptomsg))}

	var buffer bytes.Buffer

	binary.Write(&buffer, binary.BigEndian, &pkt)
	binary.Write(&buffer, binary.BigEndian, &payload)
	buffer.WriteString(cryptomsg)

	sendPktToAddr(neighborhood[key].addr, buffer.Bytes())
}

// SendPubMsg send message to public channel
func SendPubMsg(message string) {
	pkt := MeshPkt{TypePubMsg, MySID}
	pubmp := PubMsgPkt{int32(len(message))}

	var buffer bytes.Buffer

	binary.Write(&buffer, binary.BigEndian, &pkt)
	binary.Write(&buffer, binary.BigEndian, &pubmp)
	buffer.WriteString(message)

	sendPkt("255.255.255.255", 8032, buffer.Bytes())
}

// SendLookupRequest send lookup request to get node name
func SendLookupRequest(peer Neighbor) {
	pkt := MeshPkt{TypeLookupReq, MySID}
	payload := LookupReqPkt{peer.identity}

	var buffer bytes.Buffer

	binary.Write(&buffer, binary.BigEndian, &pkt)
	binary.Write(&buffer, binary.BigEndian, &payload)

	sendPktToAddr(peer.addr, buffer.Bytes())
}

// ReplyAddr get *net.UDPAddr to given MeshPkt
func ReplyAddr(hdr MeshPkt) *net.UDPAddr {
	key := hex.EncodeToString(hdr.SenderID[:])

	entry := neighborhood[key]
	return entry.addr
}

// ReplyPeer get Neighbor to given MeshPkt
func ReplyPeer(hdr MeshPkt) Neighbor {
	key := hex.EncodeToString(hdr.SenderID[:])

	return neighborhood[key]
}

// SendLookupReply send replay with node name
func SendLookupReply(peer Neighbor) {
	pkt := MeshPkt{TypeLookupRep, MySID}
	payload := LookupRepPkt{MySID, int32(len(NodeName))}

	var buffer bytes.Buffer

	binary.Write(&buffer, binary.BigEndian, &pkt)
	binary.Write(&buffer, binary.BigEndian, &payload)
	buffer.WriteString(NodeName)

	sendPktToAddr(peer.addr, buffer.Bytes())
}
