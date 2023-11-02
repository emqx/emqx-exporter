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
	} else if pos := strings.IndexRune(slice[1], '.'); pos != 0 {
		return slice[1][:pos]
	}
	return slice[1]
}
