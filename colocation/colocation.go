package colocation

import (
	"github.com/assi010/gotransip/v6/ipaddress"
	"net"
)

// Colocation struct for a Colocation
type Colocation struct {
	// List of IP ranges
	IPRanges []ipaddress.IPRange `json:"ipRanges"`
	// Colocation name
	Name string `json:"name"`
}

// colocationsWrapper struct contains a list of Colocations in it,
// this is solely used for unmarshalling
type colocationsWrapper struct {
	// array of Colocations
	Colocations []Colocation `json:"colocations"`
}

// colocationWrapper struct contains a Colocation in it,
// this is solely used for unmarshalling
type colocationWrapper struct {
	// array of Colocations
	Colocation Colocation `json:"colocation"`
}

// remoteHandsRequestWrapper is used to marshal a json request for a RemoteHandsRequest
// it encapsulates the RemoteHandsRequest with a remoteHands key
type remoteHandsRequestWrapper struct {
	RemoteHands RemoteHandsRequest `json:"remoteHands"`
}

// ipAddressWrapper struct contains an IPAddress in it,
// this is solely used for unmarshalling
type ipAddressWrapper struct {
	IPAddress ipaddress.IPAddress `json:"ipAddress"`
}

// addIpRequest struct contains an IPAddress in it,
// this is solely used for marshalling
type addIPRequest struct {
	// The IP address to register to the colocation
	IPAddress net.IP `json:"ipAddress"`
	// Reverse DNS, also known as the PTR record
	ReverseDNS string `json:"reverseDns,omitempty"`
}

// RemoteHandsRequest struct for a RemoteHandsRequest
type RemoteHandsRequest struct {
	// Name of the colocation contract
	ColoName string `json:"coloName,omitempty"`
	// Name of the person that created the remote hands request
	ContactName string `json:"contactName,omitempty"`
	// Phonenumber to contact in case of questions regarding the remotehands request
	PhoneNumber string `json:"phoneNumber,omitempty"`
	// Expected duration in minutes
	ExpectedDuration int `json:"expectedDuration,omitempty"`
	// The instructions for the datacenter engineer to perform
	Instructions string `json:"instructions,omitempty"`
}
