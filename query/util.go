package query

import (
	"net"
	"sort"

	gjson "github.com/tidwall/gjson"
)

func nullToString(d gjson.Result, p string) (r string) {
	dType := d.Get(p).Type.String()
	r = ""
	if dType == "String" {
		r = d.Get(p).Str
	}
	return r
}

func sortNets(nets []*net.IPNet) (s []string) {
	sort.SliceStable(nets, func(a, b int) bool {

		apl, _ := nets[a].Mask.Size()
		bpl, _ := nets[b].Mask.Size()
		return apl < bpl
	})

	for _, n := range nets {
		s = append(s, n.String())
	}
	sort.Strings(s)
	return s
}
