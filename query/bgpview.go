package query

import GJson "github.com/tidwall/gjson"

// GetASNDetail gets the details of an ASN.
func GetASNDetail(asn string) (*AsnDetail, error) {
	data, err := bgpViewJSON("asn", asn)

	if err != nil {
		return nil, err
	}

	lg := nullToString(data, "data.looking_glass")
	ws := nullToString(data, "data.website")

	r := &AsnDetail{
		Org:          data.Get("data.description_short").Str,
		ASN:          data.Get("data.asn").Int(),
		Country:      data.Get("data.country_code").Str,
		LookingGlass: lg,
		Website:      ws,
	}
	return r, nil
}

// GetPrefixDetail gets the details of a prefix.
func GetPrefixDetail(pfx string) (*PrefixDetail, error) {
	data, err := bgpViewJSON("prefix", pfx)

	if err != nil {
		return nil, err
	}

	rawOrigins := data.Get("data.asns")

	var origins []PrefixOrigin

	rawOrigins.ForEach(func(key, value GJson.Result) bool {
		o := PrefixOrigin{ASN: value.Get("asn").Int(), Org: value.Get("description").Str}
		origins = append(origins, o)
		return true
	})

	p := &PrefixDetail{
		Org:        data.Get("data.description_short").Str,
		Name:       data.Get("data.name").Str,
		Prefix:     data.Get("data.prefix").Str,
		Allocation: data.Get("data.rir_allocation.prefix").Str,
		RIR:        data.Get("data.rir_allocation.rir_name").Str,
		Origins:    origins,
	}
	return p, nil
}

// GetIPDetail gets details of an IP address.
func GetIPDetail(ip string) (*IPDetail, error) {
	data, err := bgpViewJSON("ip", ip)

	if err != nil {
		return nil, err
	}

	prefix := data.Get("data.prefixes").Array()[0]

	i := &IPDetail{
		IP:         data.Get("data.ip").Str,
		PTR:        data.Get("data.ptr_record").Str,
		Allocation: data.Get("data.rir_allocation.prefix").Str,
		RIR:        data.Get("data.rir_allocation.rir_name").Str,
		Org:        prefix.Get("description").Str,
		Name:       prefix.Get("name").Str,
		Prefix:     prefix.Get("prefix").Str,
		ASN:        prefix.Get("asn.asn").Int(),
	}
	return i, nil
}
