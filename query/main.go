package query

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	GJson "github.com/tidwall/gjson"
)

// AsnDetail describes a bgpview.io /asn/<asn> response.
type AsnDetail struct {
	ASN          int64
	Org          string
	Country      string
	LookingGlass string
	Website      string
}

// PrefixOrigin described a bgpview.io prefix origin.
type PrefixOrigin struct {
	ASN int64
	Org string
}

// PrefixDetail describes a bgpview.io /prefix/<prefix> response.
type PrefixDetail struct {
	Org        string
	Origins    []PrefixOrigin
	Name       string
	Prefix     string
	Allocation string
	RIR        string
}

// IPDetail describes a bgpview.io /ip/<ip> response.
type IPDetail struct {
	Allocation string
	RIR        string
	IP         string
	PTR        string
	ASN        int64
	Prefix     string
	Org        string
	Name       string
}

// ASNPrefixes is a set of IPv4 & IPv6 post-parsed & sorted addresses.
type ASNPrefixes struct {
	IPv4 []string
	IPv6 []string
}

const requestTimeout int = 15

const urlBGPView = "https://api.bgpview.io/"
const urlBGPStuff = "https://bgpstuff.net/"

var httpClient *http.Client
var headers map[string]string

func createClient() *http.Client {
	jar, err := cookiejar.New(nil)

	if err != nil {
		panic(err)
	}

	client := &http.Client{Jar: jar, Transport: &http.Transport{
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}, Timeout: time.Duration(requestTimeout) * time.Second}

	return client
}

func init() {
	httpClient = createClient()
}

func request(u string, p ...string) (b string, e error) {
	reqURL, err := url.Parse(u + strings.Join(p, "/"))

	if err != nil {
		return "", err
	}

	req, _ := http.NewRequest(http.MethodGet, reqURL.String(), nil)

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	res, err := httpClient.Do(req)

	if err != nil {
		return "", fmt.Errorf("Request to %s failed:\n%s", reqURL, err.Error())
	}

	if res.StatusCode != 200 {
		return "", fmt.Errorf("Error requesting data from %s - %s", reqURL, res.Status)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return "", fmt.Errorf("Unable to parse response from %s\n%s", reqURL, err.Error())
	}

	return string(body), nil
}

func bgpViewJSON(p ...string) (r GJson.Result, e error) {

	b, err := request(urlBGPView, p...)

	r = GJson.Parse(b)

	if err != nil {
		return r, err
	}

	if r.Get("status").Str != "ok" {
		return r, fmt.Errorf(r.Get("status_message").Str)
	}

	return r, nil
}

func bgpStuffJSON(p ...string) (r GJson.Result, e error) {
	headers = make(map[string]string)
	headers["Content-Type"] = "application/json"
	b, err := request(urlBGPStuff, p...)

	r = GJson.Parse(b)

	if err != nil {
		return r, err
	}

	return r, nil
}
