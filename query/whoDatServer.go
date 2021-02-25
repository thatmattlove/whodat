package query

import GJson "github.com/tidwall/gjson"

// GetASNDetail gets the details of an ASN.
func GetASNDetail(asn string) (*AsnDetail, error) {
	data, err := whoDatJSON("asn", asn)

	if err != nil {
		return nil, err
	}

	lg := nullToString(data, "lg")
	ws := nullToString(data, "website")

	r := &AsnDetail{
		Org:          data.Get("org").Str,
		ASN:          data.Get("asn").Int(),
		Country:      data.Get("country").Str,
		LookingGlass: lg,
		Website:      ws,
	}
	return r, nil
}

// GetPrefixDetail gets the details of a prefix.
func GetPrefixDetail(pfx string) (*PrefixDetail, error) {
	data, err := whoDatJSON("prefix", pfx)

	if err != nil {
		return nil, err
	}

	rawOrigins := data.Get("origins")

	var origins []PrefixOrigin

	rawOrigins.ForEach(func(key, value GJson.Result) bool {
		o := PrefixOrigin{ASN: value.Get("asn").Int(), Org: value.Get("org").Str, Name: value.Get("name").Str}
		origins = append(origins, o)
		return true
	})

	p := &PrefixDetail{
		Org:     data.Get("org").Str,
		Name:    data.Get("name").Str,
		Prefix:  data.Get("prefix").Str,
		RIR:     data.Get("rir").Str,
		Origins: origins,
	}
	return p, nil
}

// GetIPDetail gets details of an IP address.
func GetIPDetail(ip string) (*IPDetail, error) {
	data, err := whoDatJSON("ip", ip)

	if err != nil {
		return nil, err
	}

	i := &IPDetail{
		IP:     data.Get("ip").Str,
		PTR:    data.Get("ptr").Str,
		RIR:    data.Get("rir").Str,
		Org:    data.Get("org").Str,
		Name:   data.Get("name").Str,
		Prefix: data.Get("prefix").Str,
		ASN:    data.Get("asn").Int(),
	}
	return i, nil
}
