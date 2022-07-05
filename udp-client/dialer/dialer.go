package dialer

import "net"

// dialer package provides an interface for some of the
// functions of standard "net" library which are used of UDP
// network connections.

type NetDialer interface {
	GetConn() (*net.UDPConn, error)
}

type netDialer struct {
}

func NewNetDialer() NetDialer {
	return &netDialer{}
}

// GetConn provides UDP network conn.
func (d *netDialer) GetConn() (*net.UDPConn, error) {
	udpAddr, err := net.ResolveUDPAddr("udp4", "127.0.0.1:8090")
	if err != nil {
		return nil, err
	}
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
