package userip

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
)

func GetIP(r *http.Request) (net.IP, error) {
	realIp := r.Header.Get("X-Real-Ip")
	fmt.Println(realIp)
	userIP := net.ParseIP(realIp)
	if userIP != nil {
		fmt.Println("no user ip")
		return userIP, nil
	}

	forwardedFor := r.Header.Get("X-Forwarded-For")
	fmt.Println(forwardedFor)
	ips := strings.Split(forwardedFor, ", ")
	for _, ip := range ips {
		userIP = net.ParseIP(ip)
		if userIP != nil {
			return userIP, nil
		}
	}

	return nil, errors.New("userip: no ip found")
}
