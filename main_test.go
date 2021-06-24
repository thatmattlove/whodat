package main

import (
	"testing"

	"github.com/thatmattlove/whodat/query"
)

func TestMain(t *testing.T) {

	t.Run("IP Detail", func(t *testing.T) {
		_, err := query.GetIPDetail("1.1.1.1")
		if err != nil {
			t.Errorf(err.Error())
		}
	})
	t.Run("Prefix Detail", func(t *testing.T) {
		_, err := query.GetPrefixDetail("1.1.1.0/24")
		if err != nil {
			t.Errorf(err.Error())
		}
	})
	t.Run("ASN Detail", func(t *testing.T) {
		_, err := query.GetASNDetail("14525")
		if err != nil {
			t.Errorf(err.Error())
		}
	})
	t.Run("ASN Prefixes", func(t *testing.T) {
		_, err := query.GetASNPrefixes("14525")
		if err != nil {
			t.Errorf(err.Error())
		}
	})

}
