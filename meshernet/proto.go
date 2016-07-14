package meshernet

import (
	"net"
	"time"
)

// Various packet types for mesher protocol
const (
	TypeAnnounce  = 1 << iota // Announce packet
	TypeLookupReq             // Lookup Request
	TypeLookupRep             // Lookup Reply
	TypePrivMsg               // Private Message
	TypePubMsg                // Public Message
)

// Various service type numbers
const (
	ServiceLookup  = 1 << iota // Lookup Service
	ServicePubChat             // Public Chat Service
	ServicePrivMsg             // Private Chat Msg
)

// MeshPkt Basic Mesh Packet Header
type MeshPkt struct {
	PktPayload int8
	SenderID   [32]byte
}
type AnnouncePkt struct {
	NumServices int8
	Services    [32]int8
}
type LookupReqPkt struct {
	QueryID [32]byte
}
type LookupRepPkt struct {
	SenderID [32]byte
	NameLen  int32
	//Name     string
}

// PrivMsgPkt private message packet
type PrivMsgPkt struct {
	ReceiverId [32]byte
	MsgLen     int32
}

// PubMsgPkt Public Message Packet
type PubMsgPkt struct {
	MsgLen int32
	//message []byte
	// next payload is string bytes, not null terminated
}

// Neighbor struct describing known peers
type Neighbor struct {
	identity [32]byte
	name     string
	addr     *net.UDPAddr
	lastSeen time.Time
}

var neighborhood map[string]Neighbor
