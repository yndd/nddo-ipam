package types

import (
	"net"
	"reflect"
)

func (n *NddoipamIpamIpPrefix) GetAdminState() string {
	if reflect.ValueOf(n.AdminState).IsZero() {
		return "enabale"
	}
	return *n.AdminState
}

func (n *NddoipamIpamIpPrefix) GetDescription() string {
	if reflect.ValueOf(n.Description).IsZero() {
		return ""
	}
	return *n.Description
}

func (n *NddoipamIpamIpPrefix) GetNetworkInstance() string {
	if reflect.ValueOf(n.NetworkInstance).IsZero() {
		return "default"
	}
	return *n.NetworkInstance
}

func (n *NddoipamIpamIpPrefix) GetPrefix() string {
	if reflect.ValueOf(n.Prefix).IsZero() {
		return ""
	}
	return *n.Prefix
}

func (n *NddoipamIpamIpPrefix) GetAdresses() uint64 {
	if reflect.ValueOf(n.Adresses).IsZero() {
		return 0
	}
	return *n.Adresses
}

func (n *NddoipamIpamIpPrefix) GetStatus() string {
	if reflect.ValueOf(n.Status).IsZero() {
		return ""
	}
	return *n.Status
}

func (n *NddoipamIpamIpPrefix) GetTenant() string {
	if reflect.ValueOf(n.Tenant).IsZero() {
		return ""
	}
	return *n.Tenant
}

func (n *NddoipamIpamIpPrefix) GetPool() bool {
	if reflect.ValueOf(n.Pool).IsZero() {
		return false
	}
	return *n.Pool
}

func (n *NddoipamIpamIpPrefix) SetAdminState(s string) {
	n.AdminState = &s
}

func (n *NddoipamIpamIpPrefix) SetDescription(s string) {
	n.Description = &s
}

func (n *NddoipamIpamIpPrefix) SetNetworkInstance(s string) {
	n.NetworkInstance = &s
}

func (n *NddoipamIpamIpPrefix) SetPrefix(s string) {
	n.Prefix = &s
}

func (n *NddoipamIpamIpPrefix) SetAddresses(s uint64) {
	n.Adresses = &s
}

func (n *NddoipamIpamIpPrefix) SetStatus(s string) {
	n.Status = &s
}

func (n *NddoipamIpamIpPrefix) SetTenant(s string) {
	n.Tenant = &s
}

func (n *NddoipamIpamIpPrefix) SetPool(s bool) {
	n.Pool = &s
}

func (n *NddoipamIpamIpPrefix) GetTags() []*NddoipamIpamIpPrefixTag {
	return n.Tag
}

func (n *NddoipamIpamIpPrefix) GetTag(s string) *NddoipamIpamIpPrefixTag {
	for _, t := range n.GetTags() {
		if t.GetKey() == s {
			return t
		}
	}
	return nil
}

func (n *NddoipamIpamIpPrefix) AppendTag(s *NddoipamIpamIpPrefixTag) {
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

func (n *NddoipamIpamIpPrefixTag) GetKey() string {
	if reflect.ValueOf(n.Key).IsZero() {
		return ""
	}
	return *n.Key
}

func (n *NddoipamIpamIpPrefixTag) GetValue() string {
	if reflect.ValueOf(n.Value).IsZero() {
		return ""
	}
	return *n.Value
}

func (n *NddoipamIpamIpPrefixTag) SetValue(s string) {
	n.Value = &s
}

func (n *NddoipamIpamIpPrefix) GetAggregates() []*NddoipamIpamIpPrefixAggregate {
	return n.Aggregate
}

func (n *NddoipamIpamIpPrefix) GetAggregate(t, ni, p string) *NddoipamIpamIpPrefixAggregate {
	for _, a := range n.GetAggregates() {
		if a.GetTenant() == t && a.GetNetworkInstance() == ni && a.GetPrefix() == p {
			return a
		}
	}
	return nil
}

func (n *NddoipamIpamIpPrefix) AppendAggregate(s *NddoipamIpamIpPrefixAggregate) {
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

func (n *NddoipamIpamIpPrefixAggregate) GetTenant() string {
	if reflect.ValueOf(n.Tenant).IsZero() {
		return ""
	}
	return *n.Tenant
}

func (n *NddoipamIpamIpPrefixAggregate) GetNetworkInstance() string {
	if reflect.ValueOf(n.NetworkInstance).IsZero() {
		return ""
	}
	return *n.NetworkInstance
}

func (n *NddoipamIpamIpPrefixAggregate) GetPrefix() string {
	if reflect.ValueOf(n.Prefix).IsZero() {
		return ""
	}
	return *n.Prefix
}

func (n *NddoipamIpamIpPrefix) GetIPPrefixes() []*NddoipamIpamIpPrefixIpPrefix {
	return n.IpPrefix
}

func (n *NddoipamIpamIpPrefix) GetIPPrefix(t, ni, p string) *NddoipamIpamIpPrefixIpPrefix {
	for _, a := range n.GetIPPrefixes() {
		if a.GetTenant() == t && a.GetNetworkInstance() == ni && a.GetPrefix() == p {
			return a
		}
	}
	return nil
}

func (n *NddoipamIpamIpPrefix) AppendIPPrefix(s *NddoipamIpamIpPrefixIpPrefix) {
	found := false
	for _, a := range n.IpPrefix {
		if a.GetNetworkInstance() == s.GetNetworkInstance() && a.GetPrefix() == s.GetPrefix() && a.GetTenant() == s.GetTenant() {
			found = true
		}
	}
	if !found {
		n.IpPrefix = append(n.IpPrefix, s)
	}
}

func (n *NddoipamIpamIpPrefixIpPrefix) GetTenant() string {
	if reflect.ValueOf(n.Tenant).IsZero() {
		return ""
	}
	return *n.Tenant
}

func (n *NddoipamIpamIpPrefixIpPrefix) GetNetworkInstance() string {
	if reflect.ValueOf(n.NetworkInstance).IsZero() {
		return ""
	}
	return *n.NetworkInstance
}

func (n *NddoipamIpamIpPrefixIpPrefix) GetPrefix() string {
	if reflect.ValueOf(n.Prefix).IsZero() {
		return ""
	}
	return *n.Prefix
}

func (n *NddoipamIpamIpPrefix) Contains(ip net.IP) (bool, error) {
	_, subnet, err := net.ParseCIDR(n.GetPrefix())
	if err != nil {
		return false, err
	}
	return subnet.Contains(ip), nil
}
