package utils

import (
	"net"
)

func parseIP(ip string) (net.IP, error) {
	parsedIP := net.ParseIP(ip)

	if ipv4 := parsedIP.To4(); ipv4 != nil {
		return ipv4, nil
	}
	if parsedIP != nil{
		return parsedIP, nil
	}
	return nil, &net.ParseError{Type: "IP address", Text: ip}
}

func MakeTrustIP(trustedIP string) (string, error) {
	ip, err := parseIP(trustedIP)

	if err != nil {
		return "", err
	}
	
	var mapRenderIP = map [int]func(trustIP string) string{
		net.IPv4len: func(trustIP string) string{
			return trustIP + "/32"
		},
		net.IPv6len: func(trustIP string) string{
			return trustIP + "/32"
		},
	}

	fn, isExistKey := mapRenderIP[len(ip)]

	if isExistKey != true{
		return "", &net.ParseError{Type: "IP address", Text: trustedIP}
	}

	return fn(trustedIP), nil
}