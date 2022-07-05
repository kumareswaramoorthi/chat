package listner

import "net"

// listner package provides an interface for some of the 
// functions of standard "net" library which are used of UDP 
// network connections.

type NetListner interface {
	GetConn() (*net.UDPConn, error)
}

type netListner struct {
}

//NewNetListner function returns a netListner object.
func NewNetListner() NetListner {
	return &netListner{}
}

//GetConn provides UDP network conn.
func (d *netListner) GetConn() (*net.UDPConn, error) {
	udpAddr, err := net.ResolveUDPAddr("udp4", "127.0.0.1:8090")
	if err != nil {
		return nil, err
	}
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
