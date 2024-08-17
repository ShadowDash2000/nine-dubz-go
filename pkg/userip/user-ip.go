package userip

import (
	"errors"
	"net"
	"net/http"
	"strings"
)

func GetIP(r *http.Request) (net.IP, error) {
	realIp := r.Header.Get("X-Real-Ip")
	userIP := net.ParseIP(realIp)
	if userIP != nil {
		return userIP, nil
	}

	forwardedFor := r.Header.Get("X-Forwarded-For")
	ips := strings.Split(forwardedFor, ", ")
	for _, ip := range ips {
		userIP = net.ParseIP(ip)
		if userIP != nil {
			return userIP, nil
		}
	}

	return nil, errors.New("userip: no ip found")
}
