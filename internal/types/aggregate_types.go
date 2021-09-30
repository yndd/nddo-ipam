package types

import (
	"net"
	"reflect"

	"github.com/mikioh/ipaddr"
)

func (n *NddoipamIpamAggregate) GetAdminState() string {
	if reflect.ValueOf(n.AdminState).IsZero() {
		return "enabale"
	}
	return *n.AdminState
}

func (n *NddoipamIpamAggregate) GetDescription() string {
	if reflect.ValueOf(n.Description).IsZero() {
		return ""
	}
	return *n.Description
}

func (n *NddoipamIpamAggregate) GetNetworkInstance() string {
	if reflect.ValueOf(n.NetworkInstance).IsZero() {
		return "default"
	}
	return *n.NetworkInstance
}

func (n *NddoipamIpamAggregate) GetPrefix() string {
	if reflect.ValueOf(n.Prefix).IsZero() {
		return ""
	}
	return *n.Prefix
}

func (n *NddoipamIpamAggregate) GetPrefixes() uint64 {
	if reflect.ValueOf(n.Prefixes).IsZero() {
		return 0
	}
	return *n.Prefixes
}

func (n *NddoipamIpamAggregate) GetRirName() string {
	if reflect.ValueOf(n.RirName).IsZero() {
		return ""
	}
	return *n.RirName
}

func (n *NddoipamIpamAggregate) GetStatus() string {
	if reflect.ValueOf(n.Status).IsZero() {
		return ""
	}
	return *n.Status
}

func (n *NddoipamIpamAggregate) GetTenant() string {
	if reflect.ValueOf(n.Tenant).IsZero() {
		return ""
	}
	return *n.Tenant
}

func (n *NddoipamIpamAggregate) SetAdminState(s string) {
	n.AdminState = &s
}

func (n *NddoipamIpamAggregate) SetDescription(s string) {
	n.Description = &s
}

func (n *NddoipamIpamAggregate) SetNetworkInstance(s string) {
	n.NetworkInstance = &s
}

func (n *NddoipamIpamAggregate) SetPrefix(s string) {
	n.Prefix = &s
}

func (n *NddoipamIpamAggregate) SetPrefixes(s uint64) {
	n.Prefixes = &s
}

func (n *NddoipamIpamAggregate) SetStatus(s string) {
	n.Status = &s
}

func (n *NddoipamIpamAggregate) SetTenant(s string) {
	n.Tenant = &s
}

func (n *NddoipamIpamAggregate) GetTags() []*NddoipamIpamAggregateTag {
	return n.Tag
}

func (n *NddoipamIpamAggregate) GetTag(s string) *NddoipamIpamAggregateTag {
	for _, t := range n.GetTags() {
		if t.GetKey() == s {
			return t
		}
	}
	return nil
}

func (n *NddoipamIpamAggregate) AppendTag(s *NddoipamIpamAggregateTag) {
	found := false
	for _, a := range n.Tag {
		if a.GetKey() == s.GetKey() {
			a.SetValue(s.GetValue())
			found = true
		}
	}
	if !found {
		n.Tag = append(n.Tag, s)
	}
}

func (n *NddoipamIpamAggregateTag) GetKey() string {
	if reflect.ValueOf(n.Key).IsZero() {
		return ""
	}
	return *n.Key
}

func (n *NddoipamIpamAggregateTag) GetValue() string {
	if reflect.ValueOf(n.Value).IsZero() {
		return ""
	}
	return *n.Value
}

func (n *NddoipamIpamAggregateTag) SetValue(s string) {
	n.Value = &s
}

func (n *NddoipamIpamAggregate) ContainsAddress(ipcheck string) bool {
	_, subnet, _ := net.ParseCIDR(n.GetPrefix())
	ip := net.ParseIP(ipcheck)
	return subnet.Contains(ip)
}

func (n *NddoipamIpamAggregate) ContainsPrefix(ipprefixcheck string) bool {
	_, subnet, _ := net.ParseCIDR(n.GetPrefix())

	_, subnetCheck, _ := net.ParseCIDR(ipprefixcheck)

	p := ipaddr.NewPrefix(subnetCheck)

	// check first IP
	if !subnet.Contains(p.IP) {
		return false
	}
	// check last IP
	return subnet.Contains(p.Last())
}
