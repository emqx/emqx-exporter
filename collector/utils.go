package collector

import (
	"net/netip"
	"strings"
)

func cutNodeName(nodeName string) string {
	slice := strings.Split(nodeName, "@")
	if len(slice) != 2 {
		return nodeName
	}

	if ip, err := netip.ParseAddr(slice[1]); err == nil {
		return ip.String()
	}
	return slice[1]
}
