package query

import (
	"fmt"
	"net"

	gjson "github.com/tidwall/gjson"
)

// GetASNPrefixes gets all announced prefixes of a given ASN.
func GetASNPrefixes(asn string) (*ASNPrefixes, error) {
	data, err := bgpStuffJSON("sourced", asn)

	if err != nil {
		return nil, err
	}

	var ip4 []*net.IPNet
	var ip6 []*net.IPNet

	data.Get("Response.Sourced.Prefixes").ForEach(func(key, value gjson.Result) bool {
		ip, nw, err := net.ParseCIDR(value.Str)
		if err != nil {
			panic(fmt.Errorf("Error parsing prefix: %s", err.Error()))
		}

		if ip.To4() != nil {
			ip4 = append(ip4, nw)
		} else {
			ip6 = append(ip6, nw)
		}
		return true
	})

	p := &ASNPrefixes{
		IPv4: sortNets(ip4),
		IPv6: sortNets(ip6),
	}
	return p, nil
}
