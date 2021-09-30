package types

import (
	"errors"
	"net"
	"reflect"

	"github.com/yndd/nddo-ipam/internal/iprange"
)

func (n *NddoipamIpamIpRange) GetAdminState() string {
	if reflect.ValueOf(n.AdminState).IsZero() {
		return "enabale"
	}
	return *n.AdminState
}

func (n *NddoipamIpamIpRange) GetDescription() string {
	if reflect.ValueOf(n.Description).IsZero() {
		return ""
	}
	return *n.Description
}

func (n *NddoipamIpamIpRange) GetNetworkInstance() string {
	if reflect.ValueOf(n.NetworkInstance).IsZero() {
		return "default"
	}
	return *n.NetworkInstance
}

func (n *NddoipamIpamIpRange) GetStart() string {
	if reflect.ValueOf(n.Start).IsZero() {
		return ""
	}
	return *n.Start
}

func (n *NddoipamIpamIpRange) GetEnd() string {
	if reflect.ValueOf(n.End).IsZero() {
		return ""
	}
	return *n.End
}

func (n *NddoipamIpamIpRange) GetSize() uint64 {
	if reflect.ValueOf(n.Size).IsZero() {
		return 0
	}
	return *n.Size
}

func (n *NddoipamIpamIpRange) GetStatus() string {
	if reflect.ValueOf(n.Status).IsZero() {
		return ""
	}
	return *n.Status
}

func (n *NddoipamIpamIpRange) GetTenant() string {
	if reflect.ValueOf(n.Tenant).IsZero() {
		return ""
	}
	return *n.Tenant
}

func (n *NddoipamIpamIpRange) SetAdminState(s string) {
	n.AdminState = &s
}

func (n *NddoipamIpamIpRange) SetDescription(s string) {
	n.Description = &s
}

func (n *NddoipamIpamIpRange) SetNetworkInstance(s string) {
	n.NetworkInstance = &s
}

func (n *NddoipamIpamIpRange) SetStart(s string) {
	n.Start = &s
}

func (n *NddoipamIpamIpRange) SetEnd(s string) {
	n.End = &s
}

func (n *NddoipamIpamIpRange) SetSize(s uint64) {
	n.Size = &s
}

func (n *NddoipamIpamIpRange) SetStatus(s string) {
	n.Status = &s
}

func (n *NddoipamIpamIpRange) SetTenant(s string) {
	n.Tenant = &s
}

func (n *NddoipamIpamIpRange) GetTags() []*NddoipamIpamIpRangeTag {
	return n.Tag
}

func (n *NddoipamIpamIpRange) GetTag(s string) *NddoipamIpamIpRangeTag {
	for _, t := range n.GetTags() {
		if t.GetKey() == s {
			return t
		}
	}
	return nil
}

func (n *NddoipamIpamIpRange) AppendTag(s *NddoipamIpamIpRangeTag) {
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

func (n *NddoipamIpamIpRangeTag) GetKey() string {
	if reflect.ValueOf(n.Key).IsZero() {
		return ""
	}
	return *n.Key
}

func (n *NddoipamIpamIpRangeTag) GetValue() string {
	if reflect.ValueOf(n.Value).IsZero() {
		return ""
	}
	return *n.Value
}

func (n *NddoipamIpamIpRangeTag) SetValue(s string) {
	n.Value = &s
}

func (n *NddoipamIpamIpRange) GetAggregates() []*NddoipamIpamIpRangeAggregate {
	return n.Aggregate
}

func (n *NddoipamIpamIpRange) GetAggregate(t, ni, p string) *NddoipamIpamIpRangeAggregate {
	for _, a := range n.GetAggregates() {
		if a.GetTenant() == t && a.GetNetworkInstance() == ni && a.GetPrefix() == p {
			return a
		}
	}
	return nil
}

func (n *NddoipamIpamIpRange) AppendAggregate(s *NddoipamIpamIpRangeAggregate) {
	found := false
	for _, a := range n.Aggregate {
		if a.GetNetworkInstance() == s.GetNetworkInstance() && a.GetPrefix() == s.GetPrefix() && a.GetTenant() == s.GetTenant() {
			found = true
		}
	}
	if !found {
		n.Aggregate = append(n.Aggregate, s)
	}
}

func (n *NddoipamIpamIpRangeAggregate) GetTenant() string {
	if reflect.ValueOf(n.Tenant).IsZero() {
		return ""
	}
	return *n.Tenant
}

func (n *NddoipamIpamIpRangeAggregate) GetNetworkInstance() string {
	if reflect.ValueOf(n.NetworkInstance).IsZero() {
		return ""
	}
	return *n.NetworkInstance
}

func (n *NddoipamIpamIpRangeAggregate) GetPrefix() string {
	if reflect.ValueOf(n.Prefix).IsZero() {
		return ""
	}
	return *n.Prefix
}

func (n *NddoipamIpamIpRange) GetIPPrefixes() []*NddoipamIpamIpRangeIpPrefix {
	return n.IpPrefix
}

func (n *NddoipamIpamIpRange) GetIPPrefix(t, ni, p string) *NddoipamIpamIpRangeIpPrefix {
	for _, a := range n.GetIPPrefixes() {
		if a.GetTenant() == t && a.GetNetworkInstance() == ni && a.GetPrefix() == p {
			return a
		}
	}
	return nil
}

func (n *NddoipamIpamIpRange) AppendIPPrefix(s *NddoipamIpamIpRangeIpPrefix) {
	n.IpPrefix = append(n.IpPrefix, s)
}

func (n *NddoipamIpamIpRangeIpPrefix) GetTenant() string {
	if reflect.ValueOf(n.Tenant).IsZero() {
		return ""
	}
	return *n.Tenant
}

func (n *NddoipamIpamIpRangeIpPrefix) GetNetworkInstance() string {
	if reflect.ValueOf(n.NetworkInstance).IsZero() {
		return ""
	}
	return *n.NetworkInstance
}

func (n *NddoipamIpamIpRangeIpPrefix) GetPrefix() string {
	if reflect.ValueOf(n.Prefix).IsZero() {
		return ""
	}
	return *n.Prefix
}

/*
func (n *NddoipamIpamIpRange) Contains(ipcheck string) bool {
	ipStart := net.ParseIP(n.GetStart())
	ipEnd := net.ParseIP(n.GetEnd())
	ipCheck := net.ParseIP(ipcheck)

	start16 := ipStart.To16()
	end16 := ipEnd.To16()
	check16 := ipCheck.To16()
	if start16 == nil || end16 == nil || check16 == nil {
		return false
	}

	if bytes.Compare(check16, start16) >= 0 && bytes.Compare(check16, end16) <= 0 {
		return true
	}
	return false
}
*/

func (n *NddoipamIpamIpRange) Parse() error {
	ipStart := net.ParseIP(n.GetStart())
	ipEnd := net.ParseIP(n.GetEnd())
	r := iprange.New(ipStart, ipEnd)
	if r == nil {
		return errors.New("invalid range")
	}
	n.SetSize(r.Size().Uint64())

	return nil
}

func (n *NddoipamIpamIpRange) Contains(ip string) (bool, error) {
	ipStart := net.ParseIP(n.GetStart())
	ipEnd := net.ParseIP(n.GetEnd())
	r := iprange.New(ipStart, ipEnd)
	if r == nil {
		return false, errors.New("invalid range")
	}
	return r.Contains(net.ParseIP(ip)), nil
}
