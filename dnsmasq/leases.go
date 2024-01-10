package dnsmasq

import (
	"bufio"
	"io"
	"net"
	"strconv"
	"strings"
	"time"
)

type MacAddress net.HardwareAddr

func (hwa MacAddress) String() string {
	return net.HardwareAddr(hwa).String()
}

func (hwa MacAddress) MarshalJSON() ([]byte, error) {
	return []byte("\"" + hwa.String() + "\""), nil
}

type Lease struct {
	MacAddress MacAddress `json:"macAddress"`
	ExpireTime time.Time  `json:"expireTime"`
	IPAddress  net.Addr   `json:"ipAddress"`
	Name       string     `json:"name"`
	ClientId   string     `json:"clientId"`
}

func parseLeaseLine(line string) (*Lease, error) {
	parts := strings.Split(line, " ")
	expireEpoch, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		return nil, err
	}
	expireTime := time.Unix(expireEpoch, 0)
	maddr, err := net.ParseMAC(parts[1])
	if err != nil {
		return nil, err
	}
	macaddr := MacAddress(maddr)
	ipaddr := &net.IPAddr{net.ParseIP(parts[2]), "ip"}
	name := parts[3]
	if name == "*" {
		name = ""
	}
	clientId := parts[4]
	if clientId == "*" {
		clientId = ""
	}
	return &Lease{macaddr, expireTime, ipaddr, name, clientId}, nil
}

func ReadLeases(reader io.Reader) ([]*Lease, error) {
	scanner := bufio.NewScanner(reader)

	leases := make([]*Lease, 0, 50)
	for scanner.Scan() {
		lease, err := parseLeaseLine(scanner.Text())
		if err != nil {
			return nil, err
		}
		leases = append(leases, lease)
	}
	return leases, nil
}
