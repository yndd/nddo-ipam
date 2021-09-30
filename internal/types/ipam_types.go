package types

import (
	"net"
	"sort"
)

func (n *NddoipamIpam) SortRanges() {
	// sort the range per size
	sort.SliceStable(n.IpRange, func(i, j int) bool {
		return n.IpRange[i].GetSize() < n.IpRange[j].GetSize()
	})
}

func (n *NddoipamIpam) SortPrefixes() {
	// sort the range per size
	sort.SliceStable(n.IpPrefix, func(i, j int) bool {
		_, ipNeti, _ := net.ParseCIDR(n.IpPrefix[i].GetPrefix())
		ipMaski, _ := ipNeti.Mask.Size()
		_, ipNetj, _ := net.ParseCIDR(n.IpPrefix[i].GetPrefix())
		ipMaskj, _ := ipNetj.Mask.Size()

		return ipMaski < ipMaskj
	})
}
